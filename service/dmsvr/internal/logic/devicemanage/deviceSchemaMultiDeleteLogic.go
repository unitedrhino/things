package devicemanagelogic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
