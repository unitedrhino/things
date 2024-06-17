package record

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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

func (l *DealLogic) Deal(req *types.AlarmRecordDealReq) error {
	_, err := l.svcCtx.Rule.AlarmRecordDeal(l.ctx, utils.Copy[ud.AlarmRecordDealReq](req))
	return err
}
