package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAccessMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleAccessMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAccessMultiUpdateLogic {
	return &RoleAccessMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleAccessMultiUpdateLogic) RoleAccessMultiUpdate(in *sys.RoleAccessMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewRoleAccessRepo(l.ctx).MultiUpdate(l.ctx, in.Id, in.AccessCodes)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
