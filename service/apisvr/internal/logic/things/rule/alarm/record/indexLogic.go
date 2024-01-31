package record

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.AlarmRecordIndexReq) (resp *types.AlarmRecordIndexResp, err error) {
	ret, err := l.svcCtx.Alarm.AlarmRecordIndex(l.ctx, &rule.AlarmRecordIndexReq{
		AlarmID:   req.AlarmID,
		Page:      logic.ToRulePageRpc(req.Page),
		TimeRange: logic.ToRuleTimeRangeRpc(req.TimeRange),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AlarmRecordIndex req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.AlarmRecord, 0, len(ret.List))
	for _, v := range ret.List {
		pi := &types.AlarmRecord{
			ID:          v.Id,
			AlarmID:     v.AlarmID,
			TriggerType: v.TriggerType,
			ProductID:   v.ProductID,
			DeviceName:  v.DeviceName,
			SceneName:   v.SceneName,
			SceneID:     v.SceneID,
			Level:       v.Level,
			LastAlarm:   v.LastAlarm,
			DealState:   v.DealState,
			CreatedTime: v.CreatedTime,
		}
		pis = append(pis, pi)
	}
	return &types.AlarmRecordIndexResp{
		Total: ret.Total,
		List:  pis,
		Num:   int64(len(pis)),
	}, nil
}
