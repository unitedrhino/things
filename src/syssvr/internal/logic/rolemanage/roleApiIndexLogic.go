package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiIndexLogic {
	return &RoleApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiIndexLogic) RoleApiIndex(in *sys.RoleApiIndexReq) (*sys.RoleApiIndexResp, error) {
	data := l.svcCtx.Casbin.GetFilteredPolicy(0, cast.ToString(in.Id), ctxs.GetUserCtx(l.ctx).TenantCode, in.AppCode)
	list := make([]*sys.AuthApiInfo, 0)
	total := int64(len(data))
	if total == 0 {
		return &sys.RoleApiIndexResp{Total: total, List: list}, nil
	}

	for _, v := range data {
		list = append(list, &sys.AuthApiInfo{
			Route:  v[3],
			Method: v[4],
		})
	}
	return &sys.RoleApiIndexResp{Total: total, List: list}, nil
}
