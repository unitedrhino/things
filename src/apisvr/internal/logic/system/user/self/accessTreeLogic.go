package self

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/access/info"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AccessTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccessTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessTreeLogic {
	return &AccessTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AccessTreeLogic) AccessTree() (resp *types.AccessTreeResp, err error) {
	uc := ctxs.GetUserCtx(l.ctx)
	roleID := uc.RoleID
	if roleID == 0 {
		return nil, nil
	}
	var accessCodes []string
	if !uc.IsAdmin {
		ids, err := l.svcCtx.RoleRpc.RoleAccessIndex(l.ctx, &sys.RoleAccessIndexReq{
			Id: roleID,
		})
		if err != nil {
			return nil, err
		}
		accessCodes = ids.AccessCodes
	} else {
		ret, err := l.svcCtx.TenantRpc.TenantAccessIndex(l.ctx, &sys.TenantAccessIndexReq{
			Code: uc.TenantCode,
		})
		if err != nil {
			return nil, err
		}
		accessCodes = ret.AccessCodes
	}
	if len(accessCodes) == 0 {
		return nil, nil
	}
	ret, err := l.svcCtx.AccessRpc.AccessInfoIndex(l.ctx, &sys.AccessInfoIndexReq{Codes: accessCodes})
	if err != nil {
		return nil, err
	}
	return &types.AccessTreeResp{
		List: info.ToAccessGroupInfoTypes(ret.List),
	}, nil
}
