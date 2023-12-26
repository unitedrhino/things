package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiMultiUpdateLogic {
	return &RoleApiMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiMultiUpdateLogic) RoleApiMultiUpdate(in *sys.RoleApiMultiUpdateReq) (*sys.Response, error) {
	err := relationDB.NewRoleApiRepo(l.ctx).MultiUpdate(l.ctx, in.Id, in.AppCode, in.ApiIDs)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
