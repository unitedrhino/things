package rolelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleMenuUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuUpdateLogic {
	return &RoleMenuUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleMenuUpdateLogic) RoleMenuUpdate(in *sys.RoleMenuUpdateReq) (*sys.Response, error) {
	err := l.svcCtx.RoleModel.UpdateRoleIDMenuID(in.Id, in.MenuID)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
