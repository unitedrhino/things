package apply

import (
	"context"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DealLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DealLogic {
	return &DealLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DealLogic) Deal(req *types.UserAreaApplyDealReq) error {
	_, err := l.svcCtx.DataM.UserAreaApplyDeal(l.ctx, &sys.UserAreaApplyDealReq{
		IsApprove: req.IsApprove,
		Ids:       req.IDs,
	})
	return err
}
