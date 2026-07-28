// deviceTransferSharePolicy.go 负责设备转让时已接受分享关系的事务内处理。
package devicemanagelogic

import (
	"context"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
)

// shouldPreserveDeviceTransferShares 判断转让是否保留其他成员的分享关系。
func shouldPreserveDeviceTransferShares(isCleanData int64) bool {
	return isCleanData == def.False
}

// transferShareRepo 定义设备转让需要的分享关系存储能力。
type transferShareRepo interface {
	FindByFilter(
		ctx context.Context,
		filter relationDB.UserDeviceShareFilter,
		page *stores.PageInfo,
	) ([]*relationDB.DmUserDeviceShare, error)
	DeleteByFilter(
		ctx context.Context,
		filter relationDB.UserDeviceShareFilter,
	) error
	UpdateWithField(
		ctx context.Context,
		filter relationDB.UserDeviceShareFilter,
		updates map[string]any,
	) error
}

// applyDeviceTransferSharePolicy 在设备转让事务内更新已接受分享关系。
func applyDeviceTransferSharePolicy(
	ctx context.Context,
	repo transferShareRepo,
	devs []*devices.Core,
	recipientUserID int64,
	projectID int64,
	tenantCode string,
	preserveShares bool,
) ([]*relationDB.DmUserDeviceShare, error) {
	filter := relationDB.UserDeviceShareFilter{Devices: devs}
	shares, err := repo.FindByFilter(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	if !preserveShares {
		if err = repo.DeleteByFilter(ctx, filter); err != nil {
			return nil, err
		}
		return shares, nil
	}

	recipientFilter := relationDB.UserDeviceShareFilter{
		Devices:      devs,
		SharedUserID: recipientUserID,
	}
	if err = repo.DeleteByFilter(ctx, recipientFilter); err != nil {
		return nil, err
	}
	preservedFilter := relationDB.UserDeviceShareFilter{
		Devices:       devs,
		ExcludeUserID: recipientUserID,
	}
	err = repo.UpdateWithField(ctx, preservedFilter, map[string]any{
		"project_id":  projectID,
		"tenant_code": tenantCode,
	})
	if err != nil {
		return nil, err
	}
	return shares, nil
}
