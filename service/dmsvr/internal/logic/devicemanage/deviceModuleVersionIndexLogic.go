package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceModuleVersionIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceModuleVersionIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceModuleVersionIndexLogic {
	return &DeviceModuleVersionIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceModuleVersionIndexLogic) DeviceModuleVersionIndex(in *dm.DeviceModuleVersionIndexReq) (*dm.DeviceModuleVersionIndexResp, error) {
	f := relationDB.DeviceModuleVersionFilter{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		ModuleCode: in.ModuleCode,
	}
	total, err := relationDB.NewDeviceModuleVersionRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	ret, err := relationDB.NewDeviceModuleVersionRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	return &dm.DeviceModuleVersionIndexResp{
		Total: total,
		List:  utils.CopySlice[dm.DeviceModuleVersion](ret),
	}, nil
}
