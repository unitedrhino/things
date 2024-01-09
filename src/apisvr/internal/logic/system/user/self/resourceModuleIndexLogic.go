package self

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/module/info"
	role "github.com/i-Things/things/src/syssvr/client/rolemanage"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceModuleIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceModuleIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceModuleIndexLogic {
	return &ResourceModuleIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceModuleIndexLogic) ResourceModuleIndex() (resp *types.ModuleInfoIndexResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	roleID := uc.RoleID
	if roleID == 0 {
		return nil, nil
	}
	var moduleCodes []string
	if !uc.IsAdmin {
		codes, err := l.svcCtx.RoleRpc.RoleModuleIndex(l.ctx, &role.RoleModuleIndexReq{AppCode: uc.AppCode, Id: roleID})
		if err != nil {
			return nil, err
		}
		if len(codes.ModuleCodes) == 0 {
			return nil, nil
		}
		moduleCodes = codes.ModuleCodes
	}

	ret, err := l.svcCtx.TenantRpc.TenantAppModuleIndex(l.ctx, &sys.TenantModuleIndexReq{AppCode: uc.AppCode, ModuleCodes: moduleCodes})
	if err != nil {
		return nil, err
	}

	return &types.ModuleInfoIndexResp{
		List: info.ToModuleInfosApi(ret.List),
	}, nil
}
