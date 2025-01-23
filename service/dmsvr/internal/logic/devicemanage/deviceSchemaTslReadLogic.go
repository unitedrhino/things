package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceSchemaTslReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaTslReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaTslReadLogic {
	return &DeviceSchemaTslReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceSchemaTslReadLogic) DeviceSchemaTslRead(in *dm.DeviceSchemaTslReadReq) (*dm.DeviceSchemaTslReadResp, error) {
	if in.WithProductSchema {
		s, err := l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
		if err != nil {
			return nil, err
		}
		return &dm.DeviceSchemaTslReadResp{Tsl: s.String()}, nil
	}
	db := relationDB.NewDeviceSchemaRepo(l.ctx)
	dbSchemas, err := db.FindByFilter(l.ctx, relationDB.DeviceSchemaFilter{
		ProductID: in.ProductID, DeviceName: in.DeviceName}, nil)
	if err != nil {
		return nil, err
	}
	schemaModel := relationDB.ToDeviceSchemaDo(in.ProductID, dbSchemas)
	schemaModel.ValidateWithFmt()
	return &dm.DeviceSchemaTslReadResp{Tsl: schemaModel.String()}, nil
}
