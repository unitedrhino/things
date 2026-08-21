package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoCanBindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// canBindProjectIDDetailPrefix 标记当前用户已绑定设备的所属项目。
const canBindProjectIDDetailPrefix = "projectID="

// canBindIsShareDetail 标记当前用户已通过设备分享拥有访问关系。
const canBindIsShareDetail = "isShare=true"

func NewDeviceInfoCanBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoCanBindLogic {
	return &DeviceInfoCanBindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeviceInfoCanBind 检查设备是否可绑定，并在当前用户已有访问关系时带回前端线索。
func (l *DeviceInfoCanBindLogic) DeviceInfoCanBind(in *dm.DeviceInfoCanBindReq) (*dm.Empty, error) {
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
		ProductID:  in.Device.ProductID,
		DeviceName: in.Device.DeviceName,
	})
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}
	if di == nil {
		dii, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(ctxs.WithRoot(l.ctx), relationDB.DeviceFilter{
			DeviceNames: []string{in.Device.DeviceName},
		})
		if err != nil {
			return nil, err
		}
		di, err = l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  dii.ProductID,
			DeviceName: in.Device.DeviceName,
		})
		if err != nil {
			return nil, err
		}
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
	if err != nil {
		return nil, err
	}
	//pi, err := l.svcCtx.ProductCache.GetData(l.ctx, in.Device.ProductID)
	//if err != nil {
	//	l.Error(err)
	//	return nil, err
	//}
	ownership, err := classifyDeviceBindOwnership(l.ctx, l.svcCtx.ProjectM, l.svcCtx.UserM,
		di.TenantCode, di.ProjectID, uc, int64(dpi.DefaultProjectID))
	if err != nil {
		return nil, err
	}
	if ownership == deviceBindOwnershipBlocked { //如果在其他租户下 则已经被绑定 或 在本租户下,但是不在一个项目下也不允许绑定
		ok, err := l.currentUserHasSharedDevice(di)
		if err != nil {
			return nil, err
		}
		if ok {
			return nil, currentUserSharedDeviceError()
		}
		return nil, errors.DeviceCantBound.WithMsg("设备已被其他用户绑定。如需解绑，请按照相关流程操作。")
	}
	if ownership == deviceBindOwnershipCurrentUserOwned || ownership == deviceBindOwnershipCurrentUserBound {
		return nil, currentUserBoundDeviceError(int64(di.ProjectID))
	}
	if string(di.TenantCode) == uc.TenantCode &&
		int64(di.ProjectID) == uc.ProjectID { //如果已经绑定到自己名下则不允许重复绑定
		return nil, currentUserBoundDeviceError(int64(di.ProjectID))
	}
	return &dm.Empty{}, nil
}

// currentUserBoundDeviceError 构造当前用户已绑定错误，并携带设备当前项目。
func currentUserBoundDeviceError(projectID int64) error {
	return errors.DeviceBound.WithMsg("设备已存在，请返回设备列表查看该设备").
		AddDetailf("%s%d", canBindProjectIDDetailPrefix, projectID)
}

// currentUserSharedDeviceError 构造当前用户已获分享设备的错误，只标记共享关系。
func currentUserSharedDeviceError() error {
	return errors.DeviceBound.WithMsg("设备已存在，请返回设备列表查看该设备").
		AddDetail(canBindIsShareDetail)
}

// currentUserHasSharedDevice 查询当前用户是否已通过分享拥有这个设备。
func (l *DeviceInfoCanBindLogic) currentUserHasSharedDevice(di *dm.DeviceInfo) (bool, error) {
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if di == nil || uc.UserID == 0 {
		return false, nil
	}
	_, err := relationDB.NewUserDeviceShareRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.UserDeviceShareFilter{
		ProductID:    di.ProductID,
		DeviceName:   di.DeviceName,
		SharedUserID: uc.UserID,
	})
	if errors.Cmp(err, errors.NotFind) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
