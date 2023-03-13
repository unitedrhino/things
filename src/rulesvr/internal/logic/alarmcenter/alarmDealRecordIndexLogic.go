package alarmcenterlogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmDealRecordIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmDealRecordIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmDealRecordIndexLogic {
	return &AlarmDealRecordIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AlarmDealRecordIndexLogic) AlarmDealRecordIndex(in *rule.AlarmDealRecordIndexReq) (*rule.AlarmDealRecordIndexResp, error) {
	// todo: add your logic here and delete this line

	return &rule.AlarmDealRecordIndexResp{}, nil
}
