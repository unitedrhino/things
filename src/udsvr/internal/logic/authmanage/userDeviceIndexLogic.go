package authmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/udsvr/internal/svc"
	"github.com/i-Things/things/src/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceIndexLogic {
	return &UserDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户设备列表
func (l *UserDeviceIndexLogic) UserDeviceIndex(in *ud.Empty) (*ud.Empty, error) {
	// todo: add your logic here and delete this line

	return &ud.Empty{}, nil
}
