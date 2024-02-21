package timerEvent

import (
	"context"
	"database/sql"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	rulelogic "github.com/i-Things/things/service/udsvr/internal/logic/rule"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"time"
)

func (l *TimerHandle) SceneTiming() error {
	now := time.Now()
	return l.runWithTenant(func(ctx context.Context) error {
		ctxs.GetUserCtx(ctx).AllProject = true
		defer func() {
			ctxs.GetUserCtx(ctx).AllProject = false
		}()
		db := stores.WithNoDebug(ctx, relationDB.NewSceneTriggerTimerRepo)
		//db := relationDB.NewSceneTriggerTimerRepo(ctx)
		list, err := db.FindByFilter(ctx, relationDB.SceneTriggerTimerFilter{Status: def.True,
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                                     //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(utils.GetZeroTime(now)), stores.CmpIsNull(true)), //当天未执行的
			Repeat:      stores.CmpOr(stores.CmpBinEq(int64(now.Weekday()), 1), stores.CmpEq(0)),    //当天需要执行或只需要执行一次的
		}, nil)
		if err != nil {
			return err
		}
		for _, v := range list {
			var po = v
			do := rulelogic.PoToSceneInfoDo(&po.SceneInfo)
			if !do.When.IsHit(ctx, now, nil) {
				continue
			}
			ctxs.GoNewCtx(ctx, func(ctx context.Context) { //执行任务
				po.LastRunTime = sql.NullTime{Valid: true, Time: now}
				err := db.Update(ctx, po)
				if err != nil { //如果失败了下次还可以执行
					l.Error(err)
					return
				}
			})
		}

		l.Debug(list)
		return nil
	})
}
