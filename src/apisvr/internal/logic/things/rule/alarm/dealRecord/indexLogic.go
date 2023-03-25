package dealRecord

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.AlarmDealRecordIndexReq) (resp *types.AlarmDealRecordIndexResp, err error) {
	ret, err := l.svcCtx.Alarm.AlarmDealRecordIndex(l.ctx, &rule.AlarmDealRecordIndexReq{
		TimeRange: logic.ToRuleTimeRangeRpc(req.TimeRange),
		Page:      logic.ToRulePageRpc(req.Page),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmDealRecordIndex req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.AlarmDealRecord, 0, len(ret.List))
	for _, v := range ret.List {
		pi := &types.AlarmDealRecord{
			ID:            v.Id,
			CreatedTime:   v.CreatedTime,
			AlarmRecordID: v.AlarmRecordID,
			Result:        v.Result,
			Type:          v.Type,
			AlarmTime:     v.AlarmTime,
		}
		pis = append(pis, pi)
	}
	return &types.AlarmDealRecordIndexResp{
		Total: ret.Total,
		List:  pis,
		Num:   int64(len(pis)),
	}, nil
}
