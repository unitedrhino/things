package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	_, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx,
		relationDB.DeviceFilter{ProductID: in.ProductID, DeviceNames: []string{in.DeviceName}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsgf("找不到该设备:%v", in.DeviceName)
		}
		return nil, err
	}
	s, err := l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	if err != nil {
		return nil, err
	}
	switch schema.AffordanceType(in.Type) {
	case schema.AffordanceTypeProperty:
		_, ok := s.Property[in.Identifier]
		if ok {
			return nil, errors.Duplicate
		}
	case schema.AffordanceTypeAction:
		_, ok := s.Action[in.Identifier]
		if ok {
			return nil, errors.Duplicate
		}
	case schema.AffordanceTypeEvent:
		_, ok := s.Event[in.Identifier]
		if ok {
			return nil, errors.Duplicate
		}
	}

	po := utils.Copy[relationDB.DmDeviceSchema](in)
	po.Tag = schema.TagDevice
	if err = logic.CheckAffordance(&po.DmSchemaCore, nil); err != nil {
		return nil, err
	}
	err = relationDB.NewDeviceSchemaRepo(l.ctx).Insert(l.ctx, po)
	if err != nil {
		l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, nil)
	}
	return &dm.Empty{}, err
}
