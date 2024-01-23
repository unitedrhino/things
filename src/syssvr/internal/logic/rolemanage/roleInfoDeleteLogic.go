package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleInfoDeleteLogic {
	return &RoleInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleInfoDeleteLogic) RoleInfoDelete(in *sys.WithID) (*sys.Response, error) {
	ti, err := relationDB.NewTenantInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantInfoFilter{Code: ctxs.GetUserCtx(l.ctx).TenantCode})
	if err != nil {
		return nil, err
	}
	if ti.AdminRoleID == in.Id {
		return nil, errors.Permissions.AddMsg("超级管理员不允许删除")
	}
	err = l.RiDB.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
