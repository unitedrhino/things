package userdevicelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareCreateLogic {
	return &UserDeviceShareCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分享设备
func (l *UserDeviceShareCreateLogic) UserDeviceShareCreate(in *dm.UserDeviceShareInfo) (*dm.WithID, error) {
	// todo: add your logic here and delete this line

	return &dm.WithID{}, nil
}
