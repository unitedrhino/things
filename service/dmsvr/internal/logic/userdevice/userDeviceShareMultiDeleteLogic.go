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

type UserDeviceShareMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareMultiDeleteLogic {
	return &UserDeviceShareMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消分享设备
func (l *UserDeviceShareMultiDeleteLogic) UserDeviceShareMultiDelete(in *dm.UserDeviceShareMultiDeleteReq) (*dm.Empty, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	if in.ProjectID == 0 {
		in.ProjectID = uc.ProjectID
	}
	pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(in.ProjectID)})
	if err != nil {
		return nil, err
	}
	if pi.AdminUserID != uc.UserID { //只有所有者和被分享者才有权限操作
		return nil, errors.Permissions
	}
	err = relationDB.NewUserDeviceShareRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
		ProjectID: in.ProjectID,
		IDs:       in.Ids,
	})
	return &dm.Empty{}, err
}
