package timerEvent

import (
	"context"
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

func (l *TimerHandle) SceneTiming() error {
	now := time.Now()

	db := stores.WithNoDebug(l.ctx, relationDB.NewSceneIfTriggerRepo)
	//db := relationDB.NewSceneIfTriggerRepo(ctx)
	list, err := db.FindByFilter(l.ctx, relationDB.SceneIfTriggerFilter{Status: def.True,
		Type:        scene.TriggerTypeTimer,
		ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                                  //小于等于当前时间点(需要执行的)
		LastRunTime: stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)),                 //当天未执行的
		ExecRepeat:  stores.CmpOr(stores.CmpBinEq(int64(now.Weekday()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
	}, nil)
	if err != nil {
		return err
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
				po.LastRunTime = utils.GetEndTime(now)
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
