package authlogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthorityApiIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthorityApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorityApiIndexLogic {
	return &AuthorityApiIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthorityApiIndexLogic) AuthorityApiIndex(in *sys.AuthorityApiIndexReq) (*sys.AuthorityApiIndexResp, error) {
	data := l.svcCtx.Casbin.GetFilteredPolicy(0, cast.ToString(in.RoleID))
	list := make([]*sys.AuthorityApiInfo, 0)
	total := int64(len(data))
	if total == 0 {
		return nil, errors.NotFind.AddDetail("GetFilteredPolicy error")
	}

	for _, v := range data {
		list = append(list, &sys.AuthorityApiInfo{
			Route:  v[1],
			Method: cast.ToInt64(v[2]),
		})
	}

	return &sys.AuthorityApiIndexResp{Total: total, List: list}, nil
}
