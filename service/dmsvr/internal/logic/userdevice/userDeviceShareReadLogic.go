package userdevicelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareReadLogic {
	return &UserDeviceShareReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备分享的详情
func (l *UserDeviceShareReadLogic) UserDeviceShareRead(in *dm.UserDeviceShareReadReq) (*dm.UserDeviceShareInfo, error) {
	// todo: add your logic here and delete this line

	return &dm.UserDeviceShareInfo{}, nil
}
