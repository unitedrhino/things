package timerEvent

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/internal/domain/scene"
	rulelogic "gitee.com/unitedrhino/things/service/udsvr/internal/logic/rule"
	"gitee.com/unitedrhino/things/service/udsvr/internal/repo/relationDB"
	"github.com/observerly/dusk/pkg/dusk"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

func (l *TimerHandle) DeviceTriggerCheck() error {
	now := time.Now()
	db := stores.WithNoDebug(l.ctx, relationDB.NewSceneIfTriggerRepo)

	list, err := db.FindByFilter(l.ctx, relationDB.SceneIfTriggerFilter{
		Status:           def.True,
		Type:             scene.TriggerTypeDevice,
		FirstTriggerTime: stores.CmpIsNull(false),
		StateKeepType:    scene.StateKeepTypeDuration,
		LastRunTime:      stores.CmpIsNull(true),
	}, nil)
	if err != nil {
		l.Errorf("scene err:%v", err)
		return err
	}
	for _, v := range list {
		var po = v
		if po.SceneInfo == nil {
			logx.WithContext(l.ctx).Errorf("scene trigger device not bind scene, trigger:%v", utils.Fmt(po))
			err = relationDB.NewSceneIfTriggerRepo(l.ctx).Delete(l.ctx, po.ID)
			if err != nil {
				l.Errorf("scene err:%v", err)
			}
			continue
		}
		if v.Device.FirstTriggerTime.Time.Add(time.Duration(v.Device.StateKeep.Value) * time.Second).After(now) {
			//没有到保持时间,忽略
			continue
		}
		func() {
			err := db.UpdateWithField(l.ctx, relationDB.SceneIfTriggerFilter{ID: v.ID}, map[string]any{
				"last_run_time": now,
			})
			if err != nil {
				l.Errorf("scene err:%v", err)
			}
		}()
		po.SceneInfo.Triggers = append(po.SceneInfo.Triggers, po)
		do := rulelogic.PoToSceneInfoDo(po.SceneInfo)

		ctx := ctxs.BindTenantCode(l.ctx, string(v.SceneInfo.TenantCode), int64(v.SceneInfo.ProjectID))
		if !do.When.IsHit(ctx, now, rulelogic.NewSceneCheckRepo(l.ctx, l.svcCtx, do)) {
			continue
		}
		ctxs.GoNewCtx(ctx, func(ctx context.Context) { //执行任务
			var err error
			if err != nil { //如果失败了下次还可以执行
				logx.WithContext(ctx).Errorf("scene err:%v", err)
				return
			}
			l.SceneExec(ctx, do)
		})
	}
	return nil
}

func (l *TimerHandle) SceneTimingTenMinutes() error {
	now := time.Now()
	db := stores.WithNoDebug(l.ctx, relationDB.NewSceneIfTriggerRepo)
	list, err := db.FindByFilter(l.ctx, relationDB.SceneIfTriggerFilter{
		Status: def.True,
		Type:   scene.TriggerTypeWeather,
	}, nil)
	if err != nil {
		return err
	}
	for _, v := range list {
		var po = v
		if po.SceneInfo == nil {
			logx.WithContext(l.ctx).Errorf("trigger weather not bind scene, trigger:%v", utils.Fmt(po))
			relationDB.NewSceneIfTriggerRepo(l.ctx).Delete(l.ctx, po.ID)
			continue
		}
		po.SceneInfo.Triggers = append(po.SceneInfo.Triggers, po)
		do := rulelogic.PoToSceneInfoDo(po.SceneInfo)

		if !do.If.Triggers[0].Weather.IsHit(l.ctx, rulelogic.NewSceneCheckRepo(l.ctx, l.svcCtx, do)) {
			if po.Weather.FirstTriggerTime.Valid { //如果处于触发状态,但是现在不触发了,则需要解除触发
				err := db.UpdateWithField(l.ctx, relationDB.SceneIfTriggerFilter{ID: v.ID}, map[string]any{
					"weather_first_trigger_time": nil,
					"last_run_time":              nil,
				})
				if err != nil {
					l.Error(err)
				}
			}
			continue
		}
		if po.Weather.FirstTriggerTime.Valid { //如果已经触发过,则忽略(默认边缘触发)
			continue
		}
		func() {
			err := db.UpdateWithField(l.ctx, relationDB.SceneIfTriggerFilter{ID: v.ID}, map[string]any{
				"weather_first_trigger_time": now,
				"last_run_time":              now,
			})
			if err != nil {
				l.Error(err)
			}
		}()
		ctx := ctxs.BindTenantCode(l.ctx, string(v.SceneInfo.TenantCode), int64(v.SceneInfo.ProjectID))

		if !do.When.IsHit(ctx, now, rulelogic.NewSceneCheckRepo(l.ctx, l.svcCtx, do)) {
			continue
		}
		ctxs.GoNewCtx(ctx, func(ctx context.Context) { //执行任务
			var err error
			if err != nil { //如果失败了下次还可以执行
				logx.WithContext(ctx).Error(err)
				return
			}
			l.SceneExec(ctx, do)
		})
	}
	return nil
}

func (l *TimerHandle) SceneTiming() error {
	now := time.Now()
	db := stores.WithNoDebug(l.ctx, relationDB.NewSceneIfTriggerRepo)
	var triggerF = []relationDB.SceneIfTriggerFilter{
		{Status: def.True,
			Type:        scene.TriggerTypeTimer,
			ExecType:    stores.CmpIn(scene.ExecTypeAt, scene.ExecTypeSunSet, scene.ExecTypeSunSet),
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                  //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)), //当天未执行的
			RepeatType:  scene.RepeatTypeWeek,
			ExecRepeat:  stores.CmpOr(stores.CmpBinEq(int64(now.Weekday()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		},
		{Status: def.True,
			Type:        scene.TriggerTypeTimer,
			ExecType:    stores.CmpIn(scene.ExecTypeAt, scene.ExecTypeSunSet, scene.ExecTypeSunSet),
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                  //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)), //当天未执行的
			RepeatType:  scene.RepeatTypeMount,
			ExecRepeat:  stores.CmpOr(stores.CmpBinEq(int64(now.Day()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		},
		{Status: def.True,
			ExecType:      stores.CmpEq(scene.ExecTypeLoop),
			Type:          scene.TriggerTypeTimer,
			ExecLoopStart: stores.CmpLte(utils.TimeToDaySec(now)), //
			ExecLoopEnd:   stores.CmpGte(utils.TimeToDaySec(now)),
			LastRunTime:   stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)),
			RepeatType:    scene.RepeatTypeMount,
			ExecRepeat:    stores.CmpOr(stores.CmpBinEq(int64(now.Day()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		},
		{Status: def.True,
			ExecType:      stores.CmpEq(scene.ExecTypeLoop),
			Type:          scene.TriggerTypeTimer,
			ExecLoopStart: stores.CmpLte(utils.TimeToDaySec(now)), //
			ExecLoopEnd:   stores.CmpGte(utils.TimeToDaySec(now)),
			LastRunTime:   stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)),
			RepeatType:    scene.RepeatTypeWeek,
			ExecRepeat:    stores.CmpOr(stores.CmpBinEq(int64(now.Weekday()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		},
	}
	var list []*relationDB.UdSceneIfTrigger
	for _, v := range triggerF {
		pos, err := db.FindByFilter(l.ctx, v, nil)
		if err != nil {
			l.Error(err)
			return err
		}
		list = append(list, pos...)
	}
	l.Debugf("scene sceneTrigger now:%v list:%v", now, utils.Fmt(list))
	var sceneSet sync.Map
	for _, v := range list {
		var po = v
		ctxs.GoNewCtx(l.ctx, func(ctx context.Context) { //执行任务
			startTime := time.Now()
			log := logx.WithContext(ctx)
			defer func() {
				endTime := time.Now()
				if startTime.Add(2 * time.Second).Before(endTime) { //如果超过了2秒钟,需要记录日志
					log.Slowf("sceneTrigger use:%v  po:%v", endTime.Sub(startTime), utils.Fmt(po))
				}
			}()
			f := l.LockRunning(ctx, "scene", po.ID)
			if f == nil { //有正在执行的或redis报错,直接返回,下次重试
				log.Infof("sceneTrigger other is running:%v", utils.Fmt(po))
				return
			}
			defer func() { //避免二次释放,数据库执行完之后也可以释放,无需等执行完
				if f != nil {
					f()
					f = nil
				}
			}()
			po.SceneInfo.Triggers = append(po.SceneInfo.Triggers, po)
			do := rulelogic.PoToSceneInfoDo(po.SceneInfo)
			if po.SceneInfo == nil {
				log.Errorf("scene trigger timer not bind scene, trigger:%v", utils.Fmt(po))
				relationDB.NewSceneIfTriggerRepo(ctx).Delete(ctx, po.ID)
				return
			}
			ctx = ctxs.BindTenantCode(ctx, string(v.SceneInfo.TenantCode), 0)
			if !do.When.IsHit(ctx, now, rulelogic.NewSceneCheckRepo(l.ctx, l.svcCtx, do)) {
				return
			}

			var err error
			func() {
				defer func() { //避免二次释放,数据库执行完之后也可以释放,无需等执行完
					if f != nil {
						f()
						f = nil
					}
				}()
				switch po.Timer.ExecType {
				case scene.ExecTypeAt:
					po.LastRunTime = sql.NullTime{
						Time:  utils.GetEndTime(now),
						Valid: true,
					}
				case scene.ExecTypeLoop:
					lastRun := now.Add(time.Duration(po.Timer.ExecLoop) * time.Second)
					if utils.TimeToDaySec(lastRun) > po.Timer.ExecLoopEnd { //如果下次执行时间已经超过了结束时间,那么就到下一天开始执行
						po.LastRunTime = sql.NullTime{
							Time:  utils.GetEndTime(now),
							Valid: true,
						}
					} else { //如果当天还需要执行,则更新为下次执行时间点
						po.LastRunTime = sql.NullTime{
							Time:  lastRun,
							Valid: true,
						}
					}
				case scene.ExecTypeSunSet, scene.ExecTypeSunRises:
					po.LastRunTime = sql.NullTime{
						Time:  utils.GetEndTime(now),
						Valid: true,
					}
					func() { //更新太阳升起落下触发时间
						pi, err := l.svcCtx.ProjectCache.GetData(ctx, int64(po.SceneInfo.ProjectID))
						if err != nil {
							log.Error(err)
							return
						}
						twilight, _, err := dusk.GetLocalCivilTwilight(time.Now(), pi.Position.Longitude, pi.Position.Latitude, 0)
						if err != nil {
							log.Error(err)
							return
						}
						switch po.Timer.ExecType {
						case scene.ExecTypeSunRises:
							po.Timer.ExecAt = utils.TimeToDaySec(twilight.Until)
						case scene.ExecTypeSunSet:
							po.Timer.ExecAt = utils.TimeToDaySec(twilight.From)
						}
						po.Timer.ExecAt += po.Timer.ExecAdd
					}()
				}
				if po.Timer.ExecRepeat == 0 { //不重复执行的只执行一次
					po.Status = def.False
				}
				err = relationDB.NewSceneIfTriggerRepo(ctx).Update(ctx, po)
				if err != nil { //如果失败了下次还可以执行
					log.Errorf("scene err:%v", err)
					return
				}
				err = relationDB.NewSceneInfoRepo(ctx).UpdateWithField(ctx, relationDB.SceneInfoFilter{IDs: []int64{po.SceneID}}, map[string]any{"last_run_time": po.LastRunTime})
				if err != nil {
					log.Errorf("scene err:%v", err)
					return
				}
			}()
			if err != nil { //如果失败了下次还可以执行
				log.Errorf("scene err:%v", err)
				return
			}
			if _, ok := sceneSet.LoadOrStore(po.SceneID, struct{}{}); ok { //多个定时触发同时触发只执行一次
				log.Infof("scene 重复触发同一个场景,跳过,场景id:%v", po.SceneID)
				return
			}
			l.SceneExec(ctx, do)
		})
	}

	l.Debug(list)
	return nil
}
