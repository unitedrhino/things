package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/device/dict"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
	"strings"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	ld 	*dm.LoginDevice
	pi	*model.ProductInfo
	template dict.Template
	topics []string
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
	err = json.Unmarshal([]byte(l.pi.Template),&l.template)
	if err != nil {
		return err
	}
	//var deviceData dict.DeviceReq
	//err = json.Unmarshal([]byte(msg.Payload),&deviceData)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (l *PublishLogic) HandleThing(msg *types.Elements) error{
	if len(l.topics) < 5 || l.topics[1] != "up"{
		return errors.Parameter.AddDetail("things topic is err:"+msg.Topic)
	}
	switch l.topics[2] {
	case "property"://属性上报

	case "event"://事件上报
	case "action"://设备响应行为执行结果
	}
	return nil
}
func (l *PublishLogic) HandleOta(msg *types.Elements) error{
	return nil
}

func (l *PublishLogic) HandleDefault(msg *types.Elements) error{
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
