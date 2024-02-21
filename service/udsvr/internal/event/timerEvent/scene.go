package timerEvent

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/tools"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	rulelogic "github.com/i-Things/things/service/udsvr/internal/logic/rule"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (l *TimerHandle) SceneTiming() error {
	now := time.Now()
	return tools.RunAllTenants(l.ctx, func(ctx context.Context) error {
		ctxs.GetUserCtx(ctx).AllProject = true
		defer func() {
			ctxs.GetUserCtx(ctx).AllProject = false
		}()
		db := stores.WithNoDebug(ctx, relationDB.NewSceneTriggerTimerRepo)
		//db := relationDB.NewSceneTriggerTimerRepo(ctx)
		list, err := db.FindByFilter(ctx, relationDB.SceneTriggerTimerFilter{Status: def.True,
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                                  //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(now), stores.CmpIsNull(true)),                 //当天未执行的
			Repeat:      stores.CmpOr(stores.CmpBinEq(int64(now.Weekday()), 1), stores.CmpEq(0)), //当天需要执行或只需要执行一次的
		}, nil)
		if err != nil {
			return err
		}
		for _, v := range list {
			var po = v
			do := rulelogic.PoToSceneInfoDo(po.SceneInfo)
			if po.SceneInfo == nil {
				logx.WithContext(l.ctx).Errorf("trigger timer not bind scene, trigger:%v", utils.Fmt(po))
				continue
			}
			if !do.When.IsHit(ctx, now, nil) {
				continue
			}
			ctxs.GoNewCtx(ctx, func(ctx context.Context) { //执行任务
				po.LastRunTime = utils.GetEndTime(now)
				if po.ExecRepeat == 0 { //不重复执行的只执行一次
					po.Status = def.False
				}
				err := db.Update(ctx, po)
				if err != nil { //如果失败了下次还可以执行
					l.Error(err)
					return
				}
				l.SceneExec(ctx, do)
			})
		}

		l.Debug(list)
		return nil
	})
}

func (l *TimerHandle) SceneExec(ctx context.Context, do *scene.Info) error {
	err := do.Then.Execute(ctx, scene.ActionRepo{
		DeviceInteract: l.svcCtx.DeviceInteract,
		DeviceM:        l.svcCtx.DeviceM,
		DeviceG:        l.svcCtx.DeviceG,
		Scene:          do,
	})
	return err
}
