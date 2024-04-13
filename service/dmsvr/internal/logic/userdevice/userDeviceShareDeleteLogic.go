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
	if uds.SharedUserID != uc.UserID { //如果不是分享的对象取消分享,需要判断是否是设备的所有者,只有设备的所有者才有权取消分享
		pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(uds.ProjectID)})
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
