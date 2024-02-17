package userdevicelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareDeleteLogic {
	return &UserDeviceShareDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消分享设备
func (l *UserDeviceShareDeleteLogic) UserDeviceShareDelete(in *dm.WithID) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
