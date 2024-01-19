package self

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreaApplyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaApplyCreateLogic {
	return &AreaApplyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaApplyCreateLogic) AreaApplyCreate(req *types.UserAreaApplyCreateReq) error {
	_, err := l.svcCtx.UserRpc.UserAreaApplyCreate(l.ctx, &sys.UserAreaApplyCreateReq{
		AreaID:   req.AreaID,
		AuthType: req.AuthType,
	})
	return err
}
