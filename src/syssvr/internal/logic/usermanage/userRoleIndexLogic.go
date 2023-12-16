package usermanagelogic

import (
	"context"
	rolemanagelogic "github.com/i-Things/things/src/syssvr/internal/logic/rolemanage"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRoleIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRoleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRoleIndexLogic {
	return &UserRoleIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserRoleIndexLogic) UserRoleIndex(in *sys.UserRoleIndexReq) (*sys.UserRoleIndexResp, error) {
	ur, err := relationDB.NewUserRoleRepo(l.ctx).FindByFilter(l.ctx, relationDB.UserRoleFilter{UserID: in.UserID}, nil)
	if err != nil {
		return nil, err
	}
	var roleIDs []int64
	for _, v := range ur {
		roleIDs = append(roleIDs, v.RoleID)
	}
	rs, err := relationDB.NewRoleInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.RoleInfoFilter{IDs: roleIDs}, nil)
	if err != nil {
		return nil, err
	}
	return &sys.UserRoleIndexResp{List: rolemanagelogic.ToRoleInfosRpc(rs), Total: int64(len(roleIDs))}, nil
}
