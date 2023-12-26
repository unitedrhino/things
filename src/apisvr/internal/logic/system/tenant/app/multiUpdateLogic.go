package app

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
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

func (l *MultiUpdateLogic) MultiUpdate(req *types.TenantAppMultiUpdateReq) error {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return err
	}
	//_, err := l.svcCtx.TenantRpc.TenantAppMultiUpdate(l.ctx, &sys.TenantAppMultiUpdateReq{Code: req.Code, AppCodes: req.AppCodes})
	//return err
	return nil
}
