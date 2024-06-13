package logic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
)

func Auth(ctx context.Context, svcCtx *svc.ServiceContext, in devices.Core) (devices.Auth, error) { //鉴权
	uc := ctxs.GetUserCtx(ctx)
	if uc == nil || uc.IsSuperAdmin { //为空只可能是外部rpc调用
		return devices.AuthAll, nil
	}
	di, err := svcCtx.DeviceCache.GetData(ctx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	})
	if err != nil {
		return devices.AuthNone, err
	}
	if di.TenantCode != uc.TenantCode {
		return devices.AuthNone, errors.Permissions.AddDetail("租户号不一致")
	}
	pa := uc.ProjectAuth[di.ProjectID]
	if pa == nil { //有可能是 分享的设备及管理员
		if uc.IsAdmin {
			if di.UserID == 0 {
				return devices.AuthAll, nil
			}
			//管理员操作被绑定的设备只能操作系统功能
			return devices.AuthSystem, nil
		}
	}
	return devices.AuthAll, nil
}
