package logic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/domain/userShared"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"strings"
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

func SchemaAccess(ctx context.Context, svcCtx *svc.ServiceContext, authType def.AuthType, dev devices.Core, param map[string]any) (outParam map[string]any, err error) {
	uc := ctxs.GetUserCtx(ctx)
	if uc != nil && !uc.IsAdmin {
		di, err := svcCtx.DeviceCache.GetData(ctx, dev)
		if err != nil {
			return nil, err
		}
		_, ok := uc.ProjectAuth[di.ProjectID]
		if !ok {
			uds, err := svcCtx.UserDeviceShare.GetData(ctx, userShared.UserShareKey{
				ProductID:    dev.ProductID,
				DeviceName:   dev.DeviceName,
				SharedUserID: uc.UserID,
			})
			if err != nil {
				return nil, errors.Permissions.AddDetail(err)
			}
			if uds.AuthType == def.AuthAdmin {
				return param, nil
			}
			for k := range param {
				sp := uds.SchemaPerm[k]
				if sp == nil && strings.Contains(k, ".") { //数组类型
					k, _, _ := schema.GetArray(k)
					sp = uds.SchemaPerm[k]
				}
				if sp != nil && sp.Perm > authType {
					return nil, errors.Parameter.AddMsgf("属性:%v 没有控制权限", k)
				}
			}
			return param, nil
		}
		return param, nil
	}
	return param, nil
}
