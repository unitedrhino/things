package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareDeleteLogic {
	return &UserDeviceShareDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消分享设备
func (l *UserDeviceShareDeleteLogic) UserDeviceShareDelete(in *dm.UserDeviceShareReadReq) (*dm.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)

	f := relationDB.UserDeviceShareFilter{
		ID:         in.Id,
		DeviceName: in.Device.GetDeviceName(),
		ProductID:  in.Device.GetProductID(),
	}
	if in.Id == 0 { //如果是被分享者来获取
		f.SharedUserID = uc.UserID
	}
	uds, err := relationDB.NewUserDeviceShareRepo(l.ctx).FindOneByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	if uds.SharedUserID != uc.UserID {
		di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(ctxs.WithAllProject(l.ctx), relationDB.DeviceFilter{ProductID: uds.ProductID, DeviceNames: []string{uds.DeviceName}})
		if err != nil {
			return nil, err
		}
		if di.UserID != uc.UserID { //只有所有者和被分享者才有权限操作
			return nil, errors.Permissions
		}
	}
	err = relationDB.NewUserDeviceShareRepo(l.ctx).Delete(l.ctx, uds.ID)
	return &dm.Empty{}, err
}
