package timerEvent

import (
	"context"
	"fmt"
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
func (l *TimerHandle) LockRunning(ctx context.Context, Type string /*scene deviceTimer*/, triggerID int64) (deferF func()) {
	key := fmt.Sprintf("things:rule:%s:trigger:%d", Type, triggerID)
	ok, err := l.svcCtx.Store.SetnxExCtx(ctx, key, time.Now().Format("2006-01-02 15:04:05.999"), 5)
	if err != nil || !ok {
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		return nil
	}
	//抢到锁了
	return func() {
		_, err := l.svcCtx.Store.DelCtx(ctx, key)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}

}
