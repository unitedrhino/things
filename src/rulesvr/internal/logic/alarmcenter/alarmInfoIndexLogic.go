package alarmcenterlogic

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/logic"
	"github.com/i-Things/things/src/rulesvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlarmInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AlarmInfoRepo
}

func NewAlarmInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlarmInfoIndexLogic {
	return &AlarmInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAlarmInfoRepo(ctx),
	}
}

func (l *AlarmInfoIndexLogic) AlarmInfoIndex(in *rule.AlarmInfoIndexReq) (*rule.AlarmInfoIndexResp, error) {
	var (
		info []*rule.AlarmInfo
		size int64
		err  error
	)
	filter := relationDB.AlarmInfoFilter{
		Name: in.Name, SceneID: in.SceneID, AlarmIDs: in.AlarmIDs}
	size, err = l.AiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.AiDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}
	info = make([]*rule.AlarmInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToAlarmInfo(v))
	}
	return &rule.AlarmInfoIndexResp{
		List:  info,
		Total: size,
	}, nil
}
