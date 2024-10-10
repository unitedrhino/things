package record

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

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
