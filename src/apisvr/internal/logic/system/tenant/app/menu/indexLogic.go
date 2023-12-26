package menu

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.TenantAppMenuIndexReq) (resp *types.TenantAppMenuIndexResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppMenuIndex(l.ctx, &sys.TenantAppMenuIndexReq{
		AppCode:    req.AppCode,
		Code:       req.Code,
		ModuleCode: req.ModuleCode,
		IsRetTree:  true,
	})

	return &types.TenantAppMenuIndexResp{
		List: system.ToTenantAppMenusApi(ret.List),
	}, nil
}
