package self

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/user/info"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.UserSelfReadReq) (resp *types.UserInfo, err error) {
	var uc = ctxs.GetUserCtx(l.ctx)
	return info.NewReadLogic(l.ctx, l.svcCtx).Read(&types.UserInfoReadReq{
		UserID:     uc.UserID,
		WithRoles:  req.WithRoles,
		WithTenant: req.WithTenant,
	})
}
