package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/devices"

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
	s, err := l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	})
	if err != nil {
		return nil, err
	}
	return &dm.DeviceSchemaTslReadResp{Tsl: s.String()}, nil
}
