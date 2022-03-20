package event

import (
	"context"
	"github.com/i-Things/things/src/ddsvr/ddDef"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type InnerSubServer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInnerSubServer(svcCtx *svc.ServiceContext, ctx context.Context) *InnerSubServer {
	return &InnerSubServer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *InnerSubServer) Publish(info *ddDef.InnerPublish) error {
	return s.svcCtx.DevLink.Publish(info.Topic, []byte(info.Payload))
}
