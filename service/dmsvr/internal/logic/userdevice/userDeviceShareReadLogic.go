package userdevicelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareReadLogic {
	return &UserDeviceShareReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备分享的详情
func (l *UserDeviceShareReadLogic) UserDeviceShareRead(in *dm.UserDeviceShareReadReq) (*dm.UserDeviceShareInfo, error) {
	uc := ctxs.GetUserCtx(l.ctx)

	f := relationDB.UserDeviceShareFilter{
		ID:         in.Id,
		DeviceName: in.Device.GetDeviceName(),
		ProductID:  in.Device.GetProductID(),
		UseBy:      in.UseBy,
	}
	if in.Id == 0 { //如果是被分享者来获取
		f.SharedUserID = uc.UserID
	}
	uds, err := relationDB.NewUserDeviceShareRepo(l.ctx).FindOneByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	if !uc.IsAdmin && uds.SharedUserID != uc.UserID {
		di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(ctxs.WithAllProject(l.ctx), relationDB.DeviceFilter{ProductID: uds.ProductID, DeviceNames: []string{uds.DeviceName}})
		if err != nil {
			return nil, err
		}
		pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(di.ProjectID)})
		if err != nil {
			return nil, err
		}
		if pi.AdminUserID != uc.UserID { //只有所有者和被分享者才有权限操作
			return nil, errors.Permissions
		}
	}

	return ToUserDeviceSharePb(uds), nil
}
