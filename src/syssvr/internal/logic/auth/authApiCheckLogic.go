package authlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthApiCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthApiCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthApiCheckLogic {
	return &AuthApiCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthApiCheckLogic) AuthApiCheck(in *sys.CheckAuthReq) (*sys.Response, error) {
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
