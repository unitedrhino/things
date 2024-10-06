package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceCollectIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceCollectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceCollectIndexLogic {
	return &UserDeviceCollectIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserDeviceCollectIndexLogic) UserDeviceCollectIndex(in *dm.Empty) (*dm.UserDeviceCollectSave, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	pos, err := relationDB.NewUserDeviceCollectRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserDeviceCollectFilter{
		UserID: uc.UserID,
	}, nil)
	if err != nil {
		return nil, err
	}
	var list []*dm.DeviceCore
	for _, v := range pos {
		list = append(list, &dm.DeviceCore{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	return &dm.UserDeviceCollectSave{Devices: list}, nil
}
