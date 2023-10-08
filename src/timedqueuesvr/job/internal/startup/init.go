package startup

import (
	"context"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/svc"
	"github.com/i-Things/things/src/timedqueuesvr/job/internal/timer"
	"time"
)

func Init(svcCtx *svc.ServiceContext) error {
	return InitTimer(svcCtx)
}

func InitTimer(svcCtx *svc.ServiceContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	//ddsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
	//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
	ctx, span := ctxs.StartSpan(ctx, "InitTimer", "")
	defer span.End()
	as := clients.NewAsynqServer(svcCtx.Config.Redis)
	utils.Go(ctx, func() {
		as.Run(timer.Timed{SvcCtx: svcCtx})
	})
	return nil
}
