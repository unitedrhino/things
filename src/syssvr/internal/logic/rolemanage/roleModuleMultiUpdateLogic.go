package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleModuleMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleModuleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleModuleMultiUpdateLogic {
	return &RoleModuleMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleModuleMultiUpdateLogic) RoleModuleMultiUpdate(in *sys.RoleModuleMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewRoleModuleRepo(l.ctx).MultiUpdate(l.ctx, in.Id, in.AppCode, in.ModuleCodes)
	return &sys.Response{}, err
}
