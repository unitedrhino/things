package dmExport

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/domain/userShared"
	"strings"
)

func SchemaAccess(ctx context.Context, dc DeviceCacheT, usc UserShareCacheT, authType def.AuthType, dev devices.Core, param map[string]any) (outParam map[string]any, err error) {
	uc := ctxs.GetUserCtx(ctx)
	if uc != nil && !uc.IsAdmin {
		di, err := dc.GetData(ctx, dev)
		if err != nil {
			return nil, err
		}
		_, ok := uc.ProjectAuth[di.ProjectID]
		if !ok {
			uds, err := usc.GetData(ctx, userShared.UserShareKey{
				ProductID:    dev.ProductID,
				DeviceName:   dev.DeviceName,
				SharedUserID: uc.UserID,
			})
			if err != nil {
				return nil, errors.Permissions
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
					return nil, errors.Permissions.AddMsgf("属性:%v 没有控制权限", k)
				}
			}
			return param, nil
		}
		return param, nil
	}
	return param, nil
}
