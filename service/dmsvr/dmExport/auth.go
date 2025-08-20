package dmExport

import (
	"context"
	"strings"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/schema"
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

func AccessPerm(ctx context.Context, dc DeviceCacheT, usc UserShareCacheT, authType def.AuthType, dev devices.Core, access string) (err error) {
	uc := ctxs.GetUserCtx(ctx)
	if uc != nil && !uc.IsAdmin {
		di, err := dc.GetData(ctx, dev)
		if err != nil {
			return err
		}
		_, ok := uc.ProjectAuth[di.ProjectID]
		if !ok {
			uds, err := usc.GetData(ctx, userShared.UserShareKey{
				ProductID:    dev.ProductID,
				DeviceName:   dev.DeviceName,
				SharedUserID: uc.UserID,
			})
			if err != nil {
				return errors.Permissions
			}
			if uds.AuthType == def.AuthAdmin {
				return nil
			}
			sp := uds.AccessPerm[access]
			if sp != nil && sp.Perm > authType {
				return errors.Permissions.AddMsgf("操作:%v 没有控制权限", access)
			}
			return nil
		}
		return nil
	}
	return nil
}
