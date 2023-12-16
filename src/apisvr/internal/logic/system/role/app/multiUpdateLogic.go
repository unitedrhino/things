package app

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiUpdateLogic {
	return &MultiUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiUpdateLogic) MultiUpdate(req *types.RoleAppMultiUpdateReq) error {
	_, err := l.svcCtx.RoleRpc.RoleAppMultiUpdate(l.ctx, &sys.RoleAppMultiUpdateReq{
		Id:       req.ID,
		AppCodes: req.AppCodes,
	})
	return err
}
