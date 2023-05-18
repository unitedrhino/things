package deviceGroup

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/middlewares"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func DeviceDeleteHandle(svcCtx *svc.ServiceContext, next middlewares.HandleFunc) middlewares.HandleFunc {
	return func(ctx context.Context, value any) {
		err := svcCtx.GroupDB.DeleteDevice(ctx, value.(*devices.Core))
		logx.WithContext(ctx).Infof("DeviceDeleteHandle value:%v err:%v", utils.Fmt(value), err)
		next(ctx, value)
	}
}
