package self

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/app/info"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResourceAppIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResourceAppIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResourceAppIndexLogic {
	return &ResourceAppIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResourceAppIndexLogic) ResourceAppIndex() (resp *types.AppInfoIndexResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	roleID := uc.RoleID
	if roleID == 0 {
		return nil, nil
	}
	var appCodes []string
	if !uc.IsAdmin {
		as, err := l.svcCtx.RoleRpc.RoleAppIndex(l.ctx, &sys.RoleAppIndexReq{Id: roleID})
		if err != nil {
			return nil, err
		}
		appCodes = as.AppCodes
		if len(appCodes) == 0 {
			return nil, nil
		}
	}

	ret, err := l.svcCtx.TenantRpc.TenantAppIndex(l.ctx, &sys.TenantAppIndexReq{AppCodes: appCodes})

	return &types.AppInfoIndexResp{
		List: info.ToAppInfosTypes(ret.List),
	}, nil
}
