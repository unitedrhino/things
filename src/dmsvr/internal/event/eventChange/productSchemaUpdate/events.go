package productSchemaUpdate

import (
	"context"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func EventsHandle(svcCtx *svc.ServiceContext) any {
	return func(ctx context.Context, productID string) {
		err := svcCtx.DataUpdate.ProductSchemaUpdate(
			ctx, &events.DeviceUpdateInfo{ProductID: productID})
		if err != nil {
			logx.WithContext(ctx).Errorf("EventsHandle productID:%v err:%v", utils.Fmt(productID), err)
		}
		logx.WithContext(ctx).Infof("EventsHandle productID:%v", productID)
	}
}
