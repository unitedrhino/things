package deviceinteractlogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
)

func Auth(ctx context.Context, svcCtx *svc.ServiceContext, dev devices.Core, param map[string]any) (outParam map[string]any, err error) {
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
				return nil, errors.Parameter
			}
			if uds.AuthType == def.AuthAdmin {
				return param, nil
			}
			for k := range param {
				sp := uds.SchemaPerm[k]
				if sp != nil && sp.Perm != def.AuthReadWrite {
					return nil, errors.Parameter.AddMsgf("属性:%v 没有控制权限", k)
				}
			}
			return param, nil
		}
		return param, nil
	}
	return param, nil
}
