package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoCanBindLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoCanBindLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoCanBindLogic {
	return &DeviceInfoCanBindLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceInfoCanBindLogic) DeviceInfoCanBind(in *dm.DeviceInfoCanBindReq) (*dm.Empty, error) {
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
		ProductID:  in.Device.ProductID,
		DeviceName: in.Device.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	dpi, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
	if err != nil {
		return nil, err
	}
	if !((di.TenantCode == def.TenantCodeDefault && di.ProjectID < 3) || int64(di.ProjectID) == uc.ProjectID || di.ProjectID == dpi.DefaultProjectID) { //如果在其他租户下 则已经被绑定 或 在本租户下,但是不在一个项目下也不允许绑定
		//只有归属于default租户和自己租户的才可以
		return nil, errors.DeviceCantBound.WithMsg("设备已被其他用户绑定。如需解绑，请按照相关流程操作。")
	}
	if string(di.TenantCode) == uc.TenantCode &&
		int64(di.ProjectID) == uc.ProjectID { //如果已经绑定到自己名下则不允许重复绑定
		return nil, errors.DeviceBound.WithMsg("设备已存在，请返回设备列表查看该设备")
	}
	return &dm.Empty{}, nil
}
