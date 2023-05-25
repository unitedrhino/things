package deviceDelete

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

func DeviceGroupHandle(svcCtx *svc.ServiceContext) func(ctx context.Context, value any) {
	return func(ctx context.Context, value any) {
		err := svcCtx.GroupDB.DeleteDevice(ctx, value.(*devices.Core))
		logx.WithContext(ctx).Infof("DeviceGroupHandle value:%v err:%v", utils.Fmt(value), err)
	}
}
