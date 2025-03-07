package share

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

func ToSharePb(in *types.UserDeviceShareInfo) *dm.UserDeviceShareInfo {
	if in == nil {
		return nil
	}
	return &dm.UserDeviceShareInfo{
		Id: in.ID,
		Device: &dm.DeviceCore{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		},
		SharedUserID: in.SharedUserID,
		AuthType:     in.AuthType,
		ProjectID:    in.ProjectID,
		ExpTime:      utils.ToRpcNullInt64(in.ExpTime),
		AccessPerm:   utils.CopyMap[dm.SharePerm](in.AccessPerm),
		SchemaPerm:   utils.CopyMap[dm.SharePerm](in.SchemaPerm),
	}
}

func ToShareTypes(in *dm.UserDeviceShareInfo, ui *sys.UserInfo) *types.UserDeviceShareInfo {
	if in == nil {
		return nil
	}
	return &types.UserDeviceShareInfo{
		ID: in.Id,
		Device: types.DeviceCore{
			ProductID:  in.Device.ProductID,
			DeviceName: in.Device.DeviceName,
		},
		AuthType:          in.AuthType,
		CreatedTime:       in.CreatedTime,
		ExpTime:           utils.ToNullInt64(in.ExpTime),
		SharedUserAccount: in.SharedUserAccount,
		SharedUserID:      in.SharedUserID,
		ProjectID:         in.ProjectID,
		User:              utils.Copy[types.UserCore](ui),
		AccessPerm:        utils.CopyMap[types.SharePerm](in.AccessPerm),
		SchemaPerm:        utils.CopyMap[types.SharePerm](in.SchemaPerm),
	}
}
func ToSharesTypes(ctx context.Context, svcCtx *svc.ServiceContext, withUser bool, in []*dm.UserDeviceShareInfo) (ret []*types.UserDeviceShareInfo) {
	for _, v := range in {
		var ui *sys.UserInfo
		var err error
		if withUser {
			ui, err = svcCtx.UserC.GetData(ctx, v.SharedUserID)
			if err != nil {
				logx.WithContext(ctx).Error(err.Error())
			}
		}
		ret = append(ret, ToShareTypes(v, ui))
	}
	return
}
func ToMuitlSharePb(in *types.UserDeviceShareMultiInfo) *dm.UserDeviceShareMultiInfo {
	if in == nil {
		return nil
	}
	var dvs []*dm.DeviceShareInfo
	for _, v := range in.Devices {
		dvs = append(dvs, &dm.DeviceShareInfo{
			DeviceName: v.DeviceName,
			ProductID:  v.ProductID,
		})
	}
	return &dm.UserDeviceShareMultiInfo{
		Devices:    dvs,
		AuthType:   in.AuthType,
		ExpTime:    in.ExpTime,
		AccessPerm: utils.CopyMap[dm.SharePerm](in.AccessPerm),
		SchemaPerm: utils.CopyMap[dm.SharePerm](in.SchemaPerm),
	}
}
func ToSharesDevices(in []*types.DeviceCore) (ret []*dm.DeviceCore) {
	for _, v := range in {
		ret = append(ret, &dm.DeviceCore{
			DeviceName: v.DeviceName,
			ProductID:  v.ProductID,
		})
	}
	return ret
}
func ToMultiShareTypes(in *dm.UserDeviceShareMultiInfo) *types.UserDeviceShareMultiIndexResp {
	if in == nil {
		return nil
	}
	var dvs []*types.DeviceShareInfo
	for _, v := range in.Devices {
		dvs = append(dvs, &types.DeviceShareInfo{
			DeviceName:  v.DeviceName,
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			DeviceAlias: v.DeviceAlias.GetValue(),
			ProductImg:  v.ProductImg,
		})
	}
	return &types.UserDeviceShareMultiIndexResp{
		Devices:     dvs,
		AuthType:    in.AuthType,
		CreatedTime: in.CreatedTime,
		ExpTime:     in.ExpTime,
		AccessPerm:  utils.CopyMap[types.SharePerm](in.AccessPerm),
		SchemaPerm:  utils.CopyMap[types.SharePerm](in.SchemaPerm),
	}
}
