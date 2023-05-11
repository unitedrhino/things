package authlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthApiIndexLogic {
	return &AuthApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthApiIndexLogic) AuthApiIndex(in *sys.AuthApiIndexReq) (*sys.AuthApiIndexResp, error) {
	data := l.svcCtx.Casbin.GetFilteredPolicy(0, cast.ToString(in.RoleID))
	list := make([]*sys.AuthApiInfo, 0)
	total := int64(len(data))
	if total == 0 {
		return nil, errors.NotFind.AddDetail("GetFilteredPolicy error")
	}

	for _, v := range data {
		list = append(list, &sys.AuthApiInfo{
			Route:  v[1],
			Method: cast.ToInt64(v[2]),
		})
	}

	return &sys.AuthApiIndexResp{Total: total, List: list}, nil
}
