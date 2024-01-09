package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleMenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	RiDB *relationDB.RoleInfoRepo
	MiDB *relationDB.MenuInfoRepo
}

func NewRoleMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleMenuIndexLogic {
	return &RoleMenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		RiDB:   relationDB.NewRoleInfoRepo(ctx),
		MiDB:   relationDB.NewMenuInfoRepo(ctx),
	}
}

func (l *RoleMenuIndexLogic) RoleMenuIndex(in *sys.RoleMenuIndexReq) (*sys.RoleMenuIndexResp, error) {
	ms, err := relationDB.NewRoleMenuRepo(l.ctx).FindByFilter(l.ctx,
		relationDB.RoleMenuFilter{RoleIDs: []int64{in.Id}, AppCode: in.AppCode, ModuleCode: in.ModuleCode}, nil)
	if err != nil {
		return nil, err
	}
	var menuIDs []int64
	for _, v := range ms {
		menuIDs = append(menuIDs, v.MenuID)
	}

	return &sys.RoleMenuIndexResp{MenuIDs: menuIDs}, nil
}
