package userdevicelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareIndexLogic {
	return &UserDeviceShareIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备分享列表(只有)
func (l *UserDeviceShareIndexLogic) UserDeviceShareIndex(in *dm.UserDeviceShareIndexReq) (*dm.UserDeviceShareIndexResp, error) {
	// todo: add your logic here and delete this line

	return &dm.UserDeviceShareIndexResp{}, nil
}
