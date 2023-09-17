package rolelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RmDB *relationDB.RoleMenuRepo
}

func NewRoleMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuUpdateLogic {
	return &RoleMenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RmDB:   relationDB.NewRoleMenuRepo(ctx),
	}
}

func (l *RoleMenuUpdateLogic) RoleMenuUpdate(in *sys.RoleMenuUpdateReq) (*sys.Response, error) {
	err := l.RmDB.MultiUpdate(l.ctx, in.Id, in.MenuID)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
