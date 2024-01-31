package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCollectDeviceIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCollectDeviceIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCollectDeviceIndexLogic {
	return &UserCollectDeviceIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCollectDeviceIndexLogic) UserCollectDeviceIndex(in *ud.Empty) (*ud.UserCollectDeviceSave, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	pos, err := relationDB.NewUserCollectDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserCollectDeviceFilter{
		UserID: uc.UserID,
	}, nil)
	if err != nil {
		return nil, err
	}
	var list []*ud.DeviceCore
	for _, v := range pos {
		list = append(list, &ud.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return &ud.UserCollectDeviceSave{Devices: list}, nil
}
