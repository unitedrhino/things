package rolelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils/cast"

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
	var checkReq [][]any
	checkReq = append(checkReq, []any{cast.ToString(in.RoleID), in.Path, in.Method})
	result, err := l.svcCtx.Casbin.BatchEnforce(checkReq)
	if err != nil {
		return nil, errors.System.AddDetail("Casbin enforce error")
	}
	for _, v := range result {
		if v {
			return &sys.Response{}, nil
		}
	}

	return nil, errors.Permissions.AddDetail("权限不足")
}
