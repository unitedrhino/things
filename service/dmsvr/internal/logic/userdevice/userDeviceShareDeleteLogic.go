package userdevicelogic

import (
	"context"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
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
func (l *UserDeviceShareDeleteLogic) UserDeviceShareDelete(in *dm.WithID) (*dm.Empty, error) {
	uds, err := relationDB.NewUserDeviceShareRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	if uds.UserID != uc.UserID {
		di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: uds.ProductID, DeviceName: uds.DeviceName})
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
	err = relationDB.NewUserDeviceShareRepo(l.ctx).Delete(l.ctx, in.Id)
	return &dm.Empty{}, err
}
