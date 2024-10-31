package devicemanagelogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceSchemaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceSchemaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceSchemaIndexLogic {
	return &DeviceSchemaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备物模型列表
func (l *DeviceSchemaIndexLogic) DeviceSchemaIndex(in *dm.DeviceSchemaIndexReq) (*dm.DeviceSchemaIndexResp, error) {
	// todo: add your logic here and delete this line

	return &dm.DeviceSchemaIndexResp{}, nil
}
