package userdevicelogic

import (
	"context"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/hashicorp/go-uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserMultiDevicesShareCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserMultiDevicesShareCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMultiDevicesShareCreateLogic {
	return &UserMultiDevicesShareCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// rpc userDeviceOtaGetVersion(UserDeviceOtaGetVersionReq)returns(userDeviceOtaGetVersionResp);
func (l *UserMultiDevicesShareCreateLogic) UserMultiDevicesShareCreate(in *dm.UserMultiDevicesShareInfo) (*dm.UserMultiDevicesShareKeyword, error) {
	// 写入caches
	shareToken, _ := uuid.GenerateUUID()
	uc := ctxs.GetUserCtx(l.ctx)
	in.UserID = uc.UserID
	//判断是否有分享的权限
	pi, err := l.svcCtx.ProjectM.ProjectInfoRead(l.ctx, &sys.ProjectWithID{ProjectID: int64(uc.ProjectID)})
	if err != nil {
		return nil, err
	}
	if pi.AdminUserID != uc.UserID {
		return nil, errors.Permissions.AddMsg("只有所有者才能分享设备")
	}
	l.svcCtx.UserMultiDeviceShare.SetData(l.ctx, shareToken, in)
	return &dm.UserMultiDevicesShareKeyword{ShareToken: shareToken}, nil
}
