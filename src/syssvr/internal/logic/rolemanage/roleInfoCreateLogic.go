package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
}

func NewRoleInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleInfoCreateLogic {
	return &RoleInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
	}
}

func (l *RoleInfoCreateLogic) RoleInfoCreate(in *sys.RoleInfo) (*sys.WithID, error) {
	po := relationDB.SysTenantRoleInfo{
		Name:   in.Name,
		Desc:   in.Desc,
		Status: in.Status,
	}
	err := l.RiDB.Insert(l.ctx, &po)
	if err != nil {
		return nil, err
	}
	return &sys.WithID{Id: po.ID}, nil
}
