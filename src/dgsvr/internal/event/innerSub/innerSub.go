package innerSub

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dgsvr/internal/domain/custom"
	"github.com/i-Things/things/src/dgsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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
	var finalPayload = info.Payload
	topic := fmt.Sprintf("%s/down/%s/%s/%s", "$"+info.Handle, info.Type, info.ProductID, info.DeviceName)

	f, err := s.svcCtx.Script.GetProtoFunc(s.ctx, info.ProductID, custom.ConvertTypeDown, info.Handle, info.Type)
	if err != nil {
		s.Errorf("%s.GetProtoFunc info:%#v err:%v", utils.FuncName(), info, err)
		return err
	}
	if f != nil { //如果写了自定义函数
		finalPayload, err = f(info.Payload)
		if err != nil {
			s.Errorf("%s.Transform info:%#v err:%v", utils.FuncName(), info, err)
			return err
		}
		s.Infof("%s.transform success before:%#v after:%#v", utils.FuncName(), info.Payload, finalPayload)
		topic = fmt.Sprintf("%s/down/%s/%s/%s/%s",
			"$"+info.Handle, info.Type, custom.CustomType, info.ProductID, info.DeviceName)
	}

	return s.svcCtx.PubDev.Publish(s.ctx, topic, finalPayload)
}
