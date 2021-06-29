package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/device"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/tal-tech/go-zero/core/logx"
	"strings"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ld       *dm.LoginDevice
	pi       *model.ProductInfo
	template *device.Template
	topics   []string
}

func NewPublishLogic(ctx context.Context, svcCtx *svc.ServiceContext) LogicHandle {
	return LogicHandle(&PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

func (l *PublishLogic) initMsg(msg *types.Elements) error {
	var err error
	l.ld, err = dm.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	l.pi,err = l.svcCtx.ProductInfo.FindOneByProductID(l.ld.ProductID)
	if err != nil {
		return err
	}
	l.template, err = device.NewTemplate([]byte(l.pi.Template))
	if err != nil {
		return err
	}
	return nil
}

func(l *PublishLogic) ErrorResp(Method,clientToken string,err error){
	respTopic := fmt.Sprintf("%s/down/%s/%s/%s",
		l.topics[0],l.topics[2],l.topics[3],l.topics[4])
	payload,_ := json.Marshal(device.DeviceResp{
		Method:      Method,
		ClientToken: clientToken}.AddStatus(err))
	l.svcCtx.Mqtt.Publish(respTopic,0,false,payload)
}

func (l *PublishLogic) HandleProperty(msg *types.Elements) error{
	l.Slowf("PublishLogic|HandleProperty")
	dreq := device.DeviceReq{}
	err := utils.Unmarshal([]byte(msg.Payload), &dreq)
	if err != nil {
		return errors.Parameter.AddDetail("things topic is err:"+msg.Topic)
	}
	switch dreq.Method {
	case device.REPORT:
		tp,err := l.template.VerifyParam(dreq.Params, device.PROPERTY)
		if err != nil {
			l.ErrorResp(device.REPORT_REPLY,dreq.ClientToken,err)
		}else if len(tp) == 0 {
			l.ErrorResp(device.REPORT_REPLY,dreq.ClientToken,errors.Parameter.AddDetail("need right param"))
		}
		l.Info(tp)
	case device.REPORT_INFO:
	case device.GET_STATUS:
	default:
		return errors.Method
	}
	return nil
}

func (l *PublishLogic) HandleEvent(msg *types.Elements) error{
	l.Slowf("PublishLogic|HandleEvent")
	return nil
}
func (l *PublishLogic) HandleAction(msg *types.Elements) error{
	l.Slowf("PublishLogic|HandleAction")
	return nil
}

func (l *PublishLogic) HandleThing(msg *types.Elements) error{
	l.Slowf("PublishLogic|HandleThing")
	if len(l.topics) < 5 || l.topics[1] != "up"{
		return errors.Parameter.AddDetail("things topic is err:"+msg.Topic)
	}
	switch l.topics[2] {
	case "property"://属性上报
		return l.HandleProperty(msg)
	case "event"://事件上报
		return l.HandleEvent(msg)
	case "action"://设备响应行为执行结果
		return l.HandleAction(msg)
	default:
		return errors.Parameter.AddDetail("things topic is err:"+msg.Topic)
	}
	return nil
}
func (l *PublishLogic) HandleOta(msg *types.Elements) error{
	l.Slowf("PublishLogic|HandleOta")
	return nil
}

func (l *PublishLogic) HandleDefault(msg *types.Elements) error{
	l.Slowf("PublishLogic|HandleDefault")
	return nil
}


func (l *PublishLogic) Handle(msg *types.Elements) error {
	l.Infof("PublishLogic|req=%+v", msg)
	err := l.initMsg(msg)
	if err != nil {
		return err
	}
	l.topics = strings.Split(msg.Topic, "/")
	if len(l.topics) > 1 {
		switch l.topics[0] {
		case "$thing":
			return l.HandleThing(msg)
		case "$ota":
			return l.HandleOta(msg)
		case l.pi.ProductID:
			return l.HandleDefault(msg)
		default:
			return errors.Parameter.AddDetail(fmt.Sprintf("not suppot topic :%s",msg.Topic))
		}
	}

	fmt.Printf("template=%+v|req=%+v\n",l.template,msg.Payload)
	return nil
}
