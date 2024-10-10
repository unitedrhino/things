package userdevicelogic

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/userShared"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDeviceShareUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDeviceShareUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDeviceShareUpdateLogic {
	return &UserDeviceShareUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新权限
func (l *UserDeviceShareUpdateLogic) UserDeviceShareUpdate(in *dm.UserDeviceShareInfo) (*dm.Empty, error) {
	uds, err := relationDB.NewUserDeviceShareRepo(l.ctx).FindOne(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	uc := ctxs.GetUserCtx(l.ctx)
	di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: uds.ProductID, DeviceNames: []string{uds.DeviceName}})
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
	uds.AuthType = in.AuthType
	uds.AccessPerm = utils.CopyMap[relationDB.SharePerm](in.AccessPerm)
	uds.SchemaPerm = utils.CopyMap[relationDB.SharePerm](in.SchemaPerm)
	if uds.AccessPerm == nil {
		uds.AccessPerm = map[string]*relationDB.SharePerm{}
	}
	if uds.SchemaPerm == nil {
		uds.SchemaPerm = map[string]*relationDB.SharePerm{}
	}
	uds.ExpTime = utils.ToNullTime2(in.ExpTime)
	if err := relationDB.NewUserDeviceShareRepo(l.ctx).Update(l.ctx, uds); err != nil {
		return nil, err
	}
	l.svcCtx.UserDeviceShare.SetData(l.ctx, userShared.UserShareKey{
		ProductID:    uds.ProductID,
		DeviceName:   uds.DeviceName,
		SharedUserID: uds.SharedUserID,
	}, nil)
	return &dm.Empty{}, nil
}
