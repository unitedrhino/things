package innerSub

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dgsvr/internal/domain/custom"
	"gitee.com/i-Things/things/service/dgsvr/internal/svc"
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
	if info.Handle == devices.Gateway && info.Type == msgGateway.TypeStatus { //网关类型如果是操作子设备需要调整topic绑定状态
		err := s.GatewayHandle(info)
		if err != nil {
			s.Error(err)
		}
	}
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
func (s *InnerSubServer) GatewayHandle(info *devices.InnerPublish) error {
	var resp msgGateway.Msg
	err := utils.Unmarshal(info.Payload, &resp)
	if err != nil {
		return errors.Parameter.AddDetailf("payload unmarshal payload:%v err:%v", string(info.Payload), err)
	}
	var topics []string
	if resp.Payload != nil && len(resp.Payload.Devices) != 0 {
		for _, v := range resp.Payload.Devices {
			topics = append(topics,
				fmt.Sprintf("$thing/down/property/%s/%s", v.ProductID, v.DeviceName),
				fmt.Sprintf("$thing/down/event/%s/%s", v.ProductID, v.DeviceName),
				fmt.Sprintf("$thing/down/action/%s/%s", v.ProductID, v.DeviceName),
				fmt.Sprintf("$ota/down/upgrade/%s/%s", v.ProductID, v.DeviceName))
		}
	}
	clientID := fmt.Sprintf("%s&%s", info.ProductID, info.DeviceName)
	if resp.Method == deviceMsg.Online {
		err = s.svcCtx.MqttClient.SetClientMutSub(s.ctx, clientID, topics, 1)
	} else {
		err = s.svcCtx.MqttClient.SetClientMutUnSub(s.ctx, clientID, topics)
	}
	return err
}
