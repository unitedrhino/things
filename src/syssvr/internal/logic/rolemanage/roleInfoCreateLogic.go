package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
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

func (l *RoleInfoCreateLogic) RoleInfoCreate(in *sys.RoleInfo) (*sys.Response, error) {
	err := l.RiDB.Insert(l.ctx, &relationDB.SysRoleInfo{
		Name:   in.Name,
		Desc:   in.Desc,
		Status: in.Status,
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	return &sys.Response{}, nil
}
