package devicemanagelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// deviceBindOwnership 表示设备绑定归属校验结果。
type deviceBindOwnership int

const (
	// deviceBindOwnershipBindable 表示设备当前归属允许被本次绑定接管。
	deviceBindOwnershipBindable deviceBindOwnership = iota
	// deviceBindOwnershipBlocked 表示设备仍被有效用户或项目占用。
	deviceBindOwnershipBlocked
	// deviceBindOwnershipCurrentUserBound 表示设备已归属当前用户的其它项目。
	deviceBindOwnershipCurrentUserBound
	// deviceBindOwnershipStale 表示设备挂在已删除项目或已删除项目管理员名下。
	deviceBindOwnershipStale
)

// bindProjectReader 读取设备所属项目信息，用于判断旧归属是否仍有效。
type bindProjectReader interface {
	ProjectInfoRead(ctx context.Context, in *sys.ProjectWithID, opts ...grpc.CallOption) (*sys.ProjectInfo, error)
}

// bindUserReader 读取项目管理员信息，用于判断项目管理员是否已注销。
type bindUserReader interface {
	UserInfoRead(ctx context.Context, in *sys.UserInfoReadReq, opts ...grpc.CallOption) (*sys.UserInfo, error)
}

// classifyDeviceBindOwnership 判断设备当前归属是否应阻止新用户绑定。
func classifyDeviceBindOwnership(
	ctx context.Context,
	projectReader bindProjectReader,
	userReader bindUserReader,
	tenantCode string,
	projectID int64,
	uc *ctxs.UserCtx,
	defaultProjectID int64,
) (deviceBindOwnership, error) {
	if (tenantCode == def.TenantCodeDefault && projectID < 3) ||
		projectID == uc.ProjectID ||
		projectID == defaultProjectID {
		return deviceBindOwnershipBindable, nil
	}

	pi, err := projectReader.ProjectInfoRead(ctxs.WithRoot(ctx), &sys.ProjectWithID{ProjectID: projectID})
	if errors.Cmp(err, errors.NotFind) {
		return deviceBindOwnershipStale, nil
	}
	if err != nil {
		return deviceBindOwnershipBlocked, err
	}
	if pi == nil {
		return deviceBindOwnershipStale, nil
	}
	if pi.AdminUserID == uc.UserID {
		return deviceBindOwnershipCurrentUserBound, nil
	}

	_, err = userReader.UserInfoRead(ctxs.WithRoot(ctx), &sys.UserInfoReadReq{UserID: pi.AdminUserID})
	if errors.Cmp(err, errors.NotFind) {
		return deviceBindOwnershipStale, nil
	}
	if err != nil {
		return deviceBindOwnershipBlocked, err
	}
	return deviceBindOwnershipBlocked, nil
}

// cleanupStaleDeviceBindArtifacts 清理已注销旧归属遗留的分享、收藏和设备画像。
func cleanupStaleDeviceBindArtifacts(ctx context.Context, dev devices.Core) error {
	return stores.GetTenantConn(ctx).Transaction(func(tx *gorm.DB) error {
		err := relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(ctx, relationDB.UserDeviceShareFilter{
			ProductID:  dev.ProductID,
			DeviceName: dev.DeviceName,
		})
		if err != nil {
			return err
		}
		err = relationDB.NewDeviceProfileRepo(tx).DeleteByFilter(ctxs.WithRoot(ctx), relationDB.DeviceProfileFilter{Device: dev})
		if err != nil {
			return err
		}
		return relationDB.NewUserDeviceCollectRepo(tx).DeleteByFilter(ctx, relationDB.UserDeviceCollectFilter{
			Cores: []*devices.Core{&dev},
		})
	})
}
