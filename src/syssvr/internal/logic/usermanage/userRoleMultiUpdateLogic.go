package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleMultiUpdateLogic {
	return &UserRoleMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleMultiUpdateLogic) UserRoleMultiUpdate(in *sys.UserRoleMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewUserRoleRepo(l.ctx).MultiUpdate(l.ctx, in.UserID, in.RoleIDs)
	return &sys.Response{}, err
}
