package devicemanagelogic

import (
	"context"

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

// 批量新增物模型,只新增没有的,已有的不处理
func (l *DeviceSchemaMultiCreateLogic) DeviceSchemaMultiCreate(in *dm.DeviceSchemaMultiCreateReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
