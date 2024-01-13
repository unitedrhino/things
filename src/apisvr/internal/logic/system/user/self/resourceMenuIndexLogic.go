package self

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceMenuIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceMenuIndexLogic {
	return &ResourceMenuIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceMenuIndexLogic) ResourceMenuIndex(req *types.UserResourceWithModuleReq) (resp *types.TenantAppMenuIndexResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	roleID := uc.RoleID
	if roleID == 0 {
		return nil, nil
	}
	var menuIDs []int64
	if !uc.IsAdmin {
		ids, err := l.svcCtx.RoleRpc.RoleMenuIndex(l.ctx, &sys.RoleMenuIndexReq{AppCode: uc.AppCode, ModuleCode: req.ModuleCode})
		if err != nil {
			return nil, err
		}
		menuIDs = ids.MenuIDs
		if len(menuIDs) == 0 {
			return nil, nil
		}
	}

	ret, err := l.svcCtx.TenantRpc.TenantAppMenuIndex(l.ctx, &sys.TenantAppMenuIndexReq{AppCode: uc.AppCode, ModuleCode: req.ModuleCode, MenuIDs: menuIDs, IsRetTree: true})
	return &types.TenantAppMenuIndexResp{List: system.ToTenantAppMenusApi(ret.List)}, nil
}
