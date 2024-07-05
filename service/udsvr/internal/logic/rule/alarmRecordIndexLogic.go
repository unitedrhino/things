package rulelogic

import (
	"context"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/udsvr/internal/logic"
	"github.com/i-Things/things/service/udsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/udsvr/internal/svc"
	"github.com/i-Things/things/service/udsvr/pb/ud"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmRecordIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmRecordIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmRecordIndexLogic {
	return &AlarmRecordIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警记录
func (l *AlarmRecordIndexLogic) AlarmRecordIndex(in *ud.AlarmRecordIndexReq) (*ud.AlarmRecordIndexResp, error) {
	f := relationDB.AlarmRecordFilter{AlarmID: in.AlarmID, DealStatuses: in.DealStatus,
		Time: logic.ToTimeRange(in.TimeRange)}
	list, err := relationDB.NewAlarmRecordRepo(l.ctx).FindByFilter(l.ctx, f, logic.ToPageInfo(in.Page).
		WithDefaultOrder(stores.OrderBy{
			Field: "createdTime",
			Sort:  stores.OrderDesc,
		}))
	if err != nil {
		return nil, err
	}
	total, err := relationDB.NewAlarmRecordRepo(l.ctx).CountByFilter(l.ctx, f)
	if err != nil {
		return nil, err
	}
	return &ud.AlarmRecordIndexResp{List: utils.CopySlice[ud.AlarmRecord](list), Total: total}, nil
}
