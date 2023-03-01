package ruleenginelogic

import (
	"context"

	"github.com/i-Things/things/src/rulesvr/internal/svc"
	"github.com/i-Things/things/src/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlowInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowInfoIndexLogic {
	return &FlowInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FlowInfoIndexLogic) FlowInfoIndex(in *rule.FlowInfoIndexReq) (*rule.FlowInfoIndexResp, error) {
	// todo: add your logic here and delete this line

	return &rule.FlowInfoIndexResp{}, nil
}
