package ruleenginelogic

import (
	"context"

	"github.com/i-Things/things/service/rulesvr/internal/svc"
	"github.com/i-Things/things/service/rulesvr/pb/rule"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlowInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlowInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlowInfoUpdateLogic {
	return &FlowInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FlowInfoUpdateLogic) FlowInfoUpdate(in *rule.FlowInfo) (*rule.Empty, error) {
	// todo: add your logic here and delete this line

	return &rule.Empty{}, nil
}
