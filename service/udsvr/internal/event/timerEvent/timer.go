package timerEvent

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/tools"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type TimerHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewTimerHandle(ctx context.Context, svcCtx *svc.ServiceContext) *TimerHandle {
	return &TimerHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *TimerHandle) DeviceTimer() error {
	now := time.Now()
	return tools.RunAllTenants(l.ctx, func(ctx context.Context) error {
		ctxs.GetUserCtx(ctx).AllProject = true
		db := stores.WithNoDebug(ctx, relationDB.NewDeviceTimerInfoRepo)
		list, err := db.FindByFilter(ctx, relationDB.DeviceTimerInfoFilter{Status: def.True,
			ExecAt:      stores.CmpLte(utils.TimeToDaySec(now)),                                     //小于等于当前时间点(需要执行的)
			LastRunTime: stores.CmpOr(stores.CmpLt(utils.GetZeroTime(now)), stores.CmpIsNull(true)), //当天未执行的
			Repeat:      stores.CmpBinEq(int64(now.Weekday()), 1),                                   //当天需要执行
		}, nil)
		if err != nil {
			return err
		}
		l.Debug(list)
		return nil
	})

}
