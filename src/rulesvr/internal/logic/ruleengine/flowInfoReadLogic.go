package ruleenginelogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlowInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowInfoReadLogic {
	return &FlowInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FlowInfoReadLogic) FlowInfoRead(in *rule.FlowInfoReadReq) (*rule.FlowInfo, error) {
	// todo: add your logic here and delete this line

	return &rule.FlowInfo{}, nil
}
