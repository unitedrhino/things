package alarmcenterlogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmDealRecordCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlarmDealRecordCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmDealRecordCreateLogic {
	return &AlarmDealRecordCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 告警处理记录
func (l *AlarmDealRecordCreateLogic) AlarmDealRecordCreate(in *rule.AlarmDealRecordCreateReq) (*rule.Response, error) {
	// todo: add your logic here and delete this line

	return &rule.Response{}, nil
}
