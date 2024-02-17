package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceCollectMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceCollectMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceCollectMultiCreateLogic {
	return &UserDeviceCollectMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户收藏的设备
func (l *UserDeviceCollectMultiCreateLogic) UserDeviceCollectMultiCreate(in *dm.UserDeviceCollectSave) (*dm.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	var ucds []*relationDB.DmUserCollectDevice
	for _, v := range in.Devices {
		ucds = append(ucds, &relationDB.DmUserCollectDevice{
			UserID:     uc.UserID,
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	err := relationDB.NewUserDeviceCollectRepo(l.ctx).MultiInsert(l.ctx, ucds)
	return &dm.Empty{}, err
}
