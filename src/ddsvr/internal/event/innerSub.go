package event

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/src/ddsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type InnerSubServer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
	ctx context.Context
}

func NewInnerSubServer(svcCtx *svc.ServiceContext, ctx context.Context) *InnerSubServer {
	return &InnerSubServer{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (s *InnerSubServer) PublishToDev(info *devices.InnerPublish) error {
	topic := fmt.Sprintf("%s/down/%s/%s/%s", "$"+info.Handle, strings.Join(info.Types, "/"), info.ProductID, info.DeviceName)
	return s.svcCtx.PubDev.Publish(s.ctx, topic, info.Payload)
}
