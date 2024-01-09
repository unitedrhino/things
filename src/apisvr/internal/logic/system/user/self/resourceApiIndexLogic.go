package self

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/tenant/app/api"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceApiIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceApiIndexLogic {
	return &ResourceApiIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceApiIndexLogic) ResourceApiIndex(req *types.UserResourceWithModuleReq) (resp *types.TenantAppApiIndexResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	roleID := uc.RoleID
	if roleID == 0 {
		return nil, nil
	}
	var apiIDs []int64
	if !uc.IsAdmin {
		ids, err := l.svcCtx.RoleRpc.RoleApiIndex(l.ctx, &sys.RoleApiIndexReq{
			Id:         roleID,
			AppCode:    uc.AppCode,
			ModuleCode: req.ModuleCode,
		})
		if err != nil {
			return nil, err
		}
		apiIDs = ids.ApiIDs
		if len(apiIDs) == 0 {
			return nil, nil
		}
	}

	ret, err := l.svcCtx.TenantRpc.TenantAppApiIndex(l.ctx, &sys.TenantAppApiIndexReq{
		AppCode:    uc.AppCode,
		ModuleCode: req.ModuleCode,
		ApiIDs:     apiIDs,
	})
	if err != nil {
		return nil, err
	}
	return &types.TenantAppApiIndexResp{
		List: api.ToTenantAppApisTypes(ret.List),
	}, nil
}
