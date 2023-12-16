package user

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMultiUpdateLogic {
	return &RoleMultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleMultiUpdateLogic) RoleMultiUpdate(req *types.UserRoleMultiUpdateReq) error {
	uc := ctxs.GetUserCtx(l.ctx)
	//这里需要判断是否是租户下的超级管理员,只有租户下的超级管理员才能修改角色
	ti, err := l.svcCtx.TenantRpc.TenantInfoRead(l.ctx, &sys.ReqWithIDCode{Code: uc.TenantCode})
	if err != nil {
		return err
	}
	if ti.AdminUserID != uc.UserID {
		return errors.Permissions.AddDetail("非超级管理员不能修改角色")
	}
	_, err = l.svcCtx.UserRpc.UserRoleMultiUpdate(l.ctx, &sys.UserRoleMultiUpdateReq{UserID: req.UserID, RoleIDs: req.RoleIDs})
	return err
}
