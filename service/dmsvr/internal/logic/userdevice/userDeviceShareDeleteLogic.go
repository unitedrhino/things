package userdevicelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	if err != nil {
		return nil, err
	}
	l.svcCtx.UserDeviceShare.SetData(l.ctx, userShared.UserShareKey{
		ProductID:    uds.ProductID,
		DeviceName:   uds.DeviceName,
		SharedUserID: uds.SharedUserID,
	}, nil)
	return &dm.Empty{}, err
}
