package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/ctxs"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCollectDeviceMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCollectDeviceMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCollectDeviceMultiCreateLogic {
	return &UserCollectDeviceMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCollectDeviceMultiCreateLogic) UserCollectDeviceMultiCreate(in *ud.UserCollectDeviceSave) (*ud.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	var ucds []*relationDB.UdUserCollectDevice
	for _, v := range in.Devices {
		ucds = append(ucds, &relationDB.UdUserCollectDevice{
			UserID:     uc.UserID,
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	err := relationDB.NewUserCollectDeviceRepo(l.ctx).MultiInsert(l.ctx, ucds)
	return &ud.Empty{}, err
}
