package timerEvent

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	rulelogic "github.com/i-Things/things/service/udsvr/internal/logic/rule"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
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
		return err
	}
	for _, v := range list {
		var po = v
		if po.SceneInfo == nil {
			logx.WithContext(l.ctx).Errorf("trigger device not bind scene, trigger:%v", utils.Fmt(po))
			relationDB.NewSceneIfTriggerRepo(l.ctx).Delete(l.ctx, po.ID)
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
				l.Error(err)
			}
		}()
		po.SceneInfo.Triggers = append(po.SceneInfo.Triggers, po)
		do := rulelogic.PoToSceneInfoDo(po.SceneInfo)

		ctx := ctxs.BindTenantCode(l.ctx, string(v.SceneInfo.TenantCode), int64(v.SceneInfo.ProjectID))
		if !do.When.IsHit(ctx, now, nil) {
			continue
		}
		ctxs.GoNewCtx(ctx, func(ctx context.Context) { //执行任务
			var err error

			if err != nil { //如果失败了下次还可以执行
				l.Error(err)
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
			ExecType:    scene.ExecTypeAt,
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                  //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)), //当天未执行的
			RepeatType:  scene.RepeatTypeWeek,
			ExecRepeat:  stores.CmpOr(stores.CmpBinEq(int64(now.Weekday()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		},
		{Status: def.True,
			Type:        scene.TriggerTypeTimer,
			ExecType:    scene.ExecTypeAt,
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                  //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)), //当天未执行的
			RepeatType:  scene.RepeatTypeMount,
			ExecRepeat:  stores.CmpOr(stores.CmpBinEq(int64(now.Day()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		},
		{Status: def.True,
			ExecType:      scene.ExecTypeLoop,
			Type:          scene.TriggerTypeTimer,
			ExecLoopStart: stores.CmpLte(utils.TimeToDaySec(now)), //
			ExecLoopEnd:   stores.CmpGte(utils.TimeToDaySec(now)),
			LastRunTime:   stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)),
			RepeatType:    scene.RepeatTypeMount,
			ExecRepeat:    stores.CmpOr(stores.CmpBinEq(int64(now.Day()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		},
		{Status: def.True,
			ExecType:      scene.ExecTypeLoop,
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
			return err
		}
		list = append(list, pos...)
	}
	var sceneSet = map[int64]struct{}{}
	for _, v := range list {
		var po = v
		po.SceneInfo.Triggers = append(po.SceneInfo.Triggers, po)
		do := rulelogic.PoToSceneInfoDo(po.SceneInfo)
		if po.SceneInfo == nil {
			logx.WithContext(l.ctx).Errorf("trigger timer not bind scene, trigger:%v", utils.Fmt(po))
			relationDB.NewSceneIfTriggerRepo(l.ctx).Delete(l.ctx, po.ID)
			continue
		}
		ctx := ctxs.BindTenantCode(l.ctx, string(v.SceneInfo.TenantCode), 0)
		if !do.When.IsHit(ctx, now, nil) {
			continue
		}

		ctxs.GoNewCtx(ctx, func(ctx context.Context) { //执行任务
			f := l.LockRunning(ctx, "scene", po.ID)
			if f == nil { //有正在执行的或redis报错,直接返回,下次重试
				return
			}
			var err error
			func() {
				defer f() //数据库执行完成后就可以释放锁了
				if po.Timer.ExecType == scene.ExecTypeAt {
					po.LastRunTime = sql.NullTime{
						Time:  utils.GetEndTime(now),
						Valid: true,
					}
				} else { //间隔时间执行
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
				}
				if po.Timer.ExecRepeat == 0 { //不重复执行的只执行一次
					po.Status = def.False
				}
				err = db.Update(ctx, po)
				if err != nil { //如果失败了下次还可以执行
					l.Error(err)
					return
				}
				stores.WithNoDebug(ctx, relationDB.NewSceneInfoRepo).UpdateWithField(ctx, relationDB.SceneInfoFilter{IDs: []int64{po.SceneID}}, map[string]any{"last_run_time": time.Now()})
			}()
			if err != nil { //如果失败了下次还可以执行
				l.Error(err)
				return
			}
			if _, ok := sceneSet[po.SceneID]; ok { //多个定时触发同时触发只执行一次
				l.Infof("重复触发同一个场景,跳过,场景id:%v", po.SceneID)
				return
			}
			sceneSet[po.SceneID] = struct{}{}
			l.SceneExec(ctx, do)
		})
	}

	l.Debug(list)
	return nil
}
