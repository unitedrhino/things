package rolemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type RoleApiAuthLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRoleApiAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleApiAuthLogic {
	return &RoleApiAuthLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RoleApiAuthLogic) RoleApiAuth(in *sys.RoleApiAuthReq) (*sys.Response, error) {
	uc := ctxs.GetUserCtx(l.ctx)
	api, err := relationDB.NewTenantAppApiRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.TenantAppApiFilter{
		AppCode:   uc.AppCode,
		Route:     in.Path,
		Method:    in.Method,
		WithRoles: true,
	})
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}

	if errors.Cmp(err, errors.NotFind) { //如果没有找到,可能是不需要鉴权的接口,需要去模块接口里查询
		api, err := relationDB.NewApiInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.ApiInfoFilter{
			Route:  in.Path,
			Method: in.Method,
		})
		if err != nil {
			return nil, err
		}
		if api.IsNeedAuth == def.True { //需要鉴权的接口,但是没有查询到,说明这个租户没有权限
			return nil, errors.Permissions.AddDetail("权限不足")
		}
		return &sys.Response{}, nil
	}
	if uc.IsAdmin { //如果是租户管理员,则有权限
		return &sys.Response{}, nil
	}
	for _, v := range api.Roles {
		if v.RoleID == in.RoleID {
			return &sys.Response{}, nil
		}
	}
	return nil, errors.Permissions.AddDetail("权限不足")
}
