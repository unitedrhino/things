package timerEvent

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/tools"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/domain/scene"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"time"
)

func (l *TimerHandle) DeviceTimer() error {
	now := time.Now()
	return tools.RunAllTenants(l.ctx, func(ctx context.Context) error {
		ctxs.GetUserCtx(ctx).AllProject = true
		defer func() {
			ctxs.GetUserCtx(ctx).AllProject = false
		}()
		db := stores.WithNoDebug(ctx, relationDB.NewDeviceTimerInfoRepo)
		list, err := db.FindByFilter(ctx, relationDB.DeviceTimerInfoFilter{Status: def.True,
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                                     //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(utils.GetZeroTime(now)), stores.CmpIsNull(true)), //当天未执行的
			Repeat:      stores.CmpBinEq(int64(now.Weekday()), 1),                                   //当天需要执行
		}, nil)
		if err != nil || len(list) == 0 {
			return err
		}
		for _, v := range list {
			po := v
			ctxs.GoNewCtx(ctx, func(ctx context.Context) {
				f := l.LockRunning(ctx, "deviceTimer", po.ID)
				if f == nil { //有正在执行的或redis报错,直接返回,下次重试
					return
				}
				var err error
				func() {
					defer f() //数据库执行完成后就可以释放锁了
					po.LastRunTime = utils.GetEndTime(now)
					if po.ExecRepeat == 0 { //不重复执行的只执行一次
						po.Status = def.False
					}
					db.Update(ctx, po)
					if err != nil { //如果失败了下次还可以执行
						l.Error(err)
						return
					}
				}()
				if err != nil { //如果失败了下次还可以执行
					l.Error(err)
					return
				}
				exec := scene.ActionDevice{
					ProductID:   po.ProductID,
					SelectType:  scene.SelectDeviceFixed,
					DeviceNames: []string{po.DeviceName},
					Type:        scene.ActionDeviceType(po.ActionType),
					DataID:      po.DataID,
					Value:       po.Value,
				}
				err = exec.Execute(ctx, scene.ActionRepo{
					DeviceInteract: l.svcCtx.DeviceInteract,
					DeviceM:        l.svcCtx.DeviceM,
					DeviceG:        l.svcCtx.DeviceG,
				})
				if err != nil {
					l.Error(err)
					return
				}
			})
		}
		l.Debug(list)
		return nil
	})
}
