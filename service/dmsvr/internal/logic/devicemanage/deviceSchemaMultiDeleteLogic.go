package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceSchemaMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaMultiDeleteLogic {
	return &DeviceSchemaMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除设备物模型
func (l *DeviceSchemaMultiDeleteLogic) DeviceSchemaMultiDelete(in *dm.DeviceSchemaMultiDeleteReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithAllProject(l.ctx)
	pos, err := relationDB.NewDeviceSchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceSchemaFilter{
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		Identifiers: in.Identifiers,
	}, nil)
	if err != nil {
		return nil, err
	}
	s := relationDB.ToDeviceSchemaDo(in.ProductID, pos)
	if len(s.Properties) > 0 {
		err = l.svcCtx.SchemaManaRepo.DeleteDeviceProperty(l.ctx, in.ProductID, in.DeviceName, s.Properties)
		if err != nil {
			return nil, err
		}
	}
	err = relationDB.NewDeviceSchemaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.DeviceSchemaFilter{
		ProductID:   in.ProductID,
		DeviceName:  in.DeviceName,
		Identifiers: in.Identifiers,
	})
	if err != nil {
		l.svcCtx.DeviceSchemaRepo.SetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, nil)
	}
	return &dm.Empty{}, err
}
