package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceCollectMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceCollectMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceCollectMultiDeleteLogic {
	return &UserDeviceCollectMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeviceCollectMultiDeleteLogic) UserDeviceCollectMultiDelete(in *dm.UserDeviceCollectSave) (*dm.Empty, error) {
	var ds []*devices.Core
	for _, v := range in.Devices {
		ds = append(ds, &devices.Core{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	uc := ctxs.GetUserCtx(l.ctx)

	err := relationDB.NewUserDeviceCollectRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserDeviceCollectFilter{UserID: uc.UserID, Cores: ds})
	return &dm.Empty{}, err
}
