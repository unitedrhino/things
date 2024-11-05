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

type DeviceSchemaMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaMultiCreateLogic {
	return &DeviceSchemaMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceSchemaMultiCreateLogic) ruleCheck(in *dm.DeviceSchemaMultiCreateReq) (ret []*relationDB.DmDeviceSchema, err error) {
	_, err = relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx,
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
	for _, v := range in.List {
		switch schema.AffordanceType(v.Type) {
		case schema.AffordanceTypeProperty:
			_, ok := s.Property[v.Identifier]
			if ok {
				continue
			}
		case schema.AffordanceTypeAction:
			_, ok := s.Action[v.Identifier]
			if ok {
				continue
			}
		case schema.AffordanceTypeEvent:
			_, ok := s.Event[v.Identifier]
			if ok {
				continue
			}
		}

		po := utils.Copy[relationDB.DmDeviceSchema](v)
		po.Tag = schema.TagDevice
		if err = logic.CheckAffordance(&po.DmSchemaCore, nil); err != nil {
			return nil, err
		}
		ret = append(ret, po)
	}
	return
}

// 批量新增物模型,只新增没有的,已有的不处理
func (l *DeviceSchemaMultiCreateLogic) DeviceSchemaMultiCreate(in *dm.DeviceSchemaMultiCreateReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithAllProject(l.ctx)
	pos, err := l.ruleCheck(in)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewDeviceSchemaRepo(l.ctx).MultiInsert(l.ctx, pos)
	if err != nil {
		l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, nil)
	}
	return &dm.Empty{}, err
}
