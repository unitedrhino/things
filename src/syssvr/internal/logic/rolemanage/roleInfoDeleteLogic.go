package rolemanagelogic

import (
	"context"
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
	err := l.RiDB.Delete(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
