package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/core/shared/ctxs"
	"gitee.com/i-Things/core/shared/devices"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCollectDeviceMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCollectDeviceMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCollectDeviceMultiDeleteLogic {
	return &UserCollectDeviceMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCollectDeviceMultiDeleteLogic) UserCollectDeviceMultiDelete(in *ud.UserCollectDeviceSave) (*ud.Empty, error) {
	var ds []*devices.Core
	for _, v := range in.Devices {
		ds = append(ds, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	uc := ctxs.GetUserCtx(l.ctx)

	err := relationDB.NewUserCollectDeviceRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserCollectDeviceFilter{UserID: uc.UserID, Cores: ds})
	return &ud.Empty{}, err
}
