package authority

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthorityApiIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthorityApiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorityApiIndexLogic {
	return &AuthorityApiIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthorityApiIndexLogic) AuthorityApiIndex() (resp *types.AuthorityApiIndexResp, err error) {
	RoleId := cast.ToString(userHeader.GetUserCtx(l.ctx).Role)
	data := l.svcCtx.Casbin.GetFilteredPolicy(0, RoleId)
	resp = &types.AuthorityApiIndexResp{}
	resp.Total = int64(len(data))
	if resp.Total > 0 {
		for _, v := range data {
			resp.List = append(resp.List, &types.AuthorityApiInfo{
				Path:   v[1],
				Method: v[2],
			})
		}
	}
	return resp, nil
}
