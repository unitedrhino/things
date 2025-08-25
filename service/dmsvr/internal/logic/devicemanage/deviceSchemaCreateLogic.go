package devicemanagelogic

import (
	"context"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceSchemaCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaCreateLogic {
	return &DeviceSchemaCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增设备
func (l *DeviceSchemaCreateLogic) DeviceSchemaCreate(in *dm.DeviceSchema) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithAllProject(l.ctx)
	pos, err := ruleCheck(l.ctx, l.svcCtx, &dm.DeviceSchemaMultiCreateReq{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		List:       []*dm.DeviceSchema{in},
	})
	if len(pos) == 0 {
		return &dm.Empty{}, err
	}
	err = relationDB.NewDeviceSchemaRepo(l.ctx).Insert(l.ctx, pos[0])
	if err != nil {
		l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, nil)
	}
	return &dm.Empty{}, err
}
