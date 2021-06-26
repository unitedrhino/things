package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dmsvr/device/dict"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"github.com/spf13/cast"
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

func CompareType(val interface{}, define *dict.Define) ( interface{} ,bool){
	switch define.Type {
	case dict.BOOL:
		switch val.(type) {
		case bool:
			return val.(bool),true
		case json.Number:
			num :=val.(json.Number).String()
			if  num == "0" {
				return false,true
			}else {
				return true,true
			}
		}
	case dict.INT:
		if num,ok:=val.(json.Number);!ok{
			return nil, false
		}else {
			ret,err := num.Int64()
			if err != nil {
				return nil, false
			}
			return ret,true
		}
	case dict.FLOAT:
		if num,ok:=val.(json.Number);!ok{
			return nil, false
		}else {
			ret,err := num.Float64()
			if err != nil {
				return nil, false
			}
			return ret,true
		}
	case dict.STRING:
		if str,ok:=val.(string);!ok{
			return nil, false
		}else {
			return str,true
		}
	case dict.ENUM://枚举类型 报文中传递的是数字
		if num,ok:=val.(json.Number);!ok{
			return nil, false
		}else {
			ret,err := num.Int64()
			if err != nil {
				return nil, false
			}
			return ret,true
		}
	case dict.TIMESTAMP:
		switch val.(type) {
		case json.Number:
			ret,err := val.(json.Number).Int64()
			if err != nil {
				return nil, false
			}
			return ret,true
		case string:
			ret,err := cast.ToInt64E(val)
			if err != nil {
				return nil, false
			}
			return ret,true
		}
	case dict.STRUCT:
		if stru,ok := val.(map[string]interface {});!ok{
			return nil, false
		}else {
			getParam := make(map[string]interface{},len(stru))
			for _,sv := range define.Specs{
				for k,v :=range stru {
					if sv.ID == k{
						param,err := CompareType(v,&sv.DataType)
						if err != false {
							getParam[k] = param
						}
					}
				}
			}
			return getParam,true
		}
	case dict.ARRAY:
		if arr,ok := val.([]interface {});!ok{
			return nil, false
		}else {
			getParam := make([]interface {},len(arr))
			for _,v :=range arr {
				param,err := CompareType(v,define.ArrayInfo)
				if err != true {
					getParam = append(getParam,param)
				}
			}
			return getParam, true
		}
	}
	return nil, false
}

func VeryfyPropertyReportReq(template *dict.Template,req *dict.DeviceReq) error{
	if len(req.Params) == 0 {
		return errors.Parameter.AddDetail("need add params")
	}
	getParam := make(map[string]interface{},len(req.Params))
	for k,v := range req.Params{
		for _,property := range template.Properties{
			if k == property.ID{
				param,err := CompareType(v,&property.Define)
				if err == false {
					getParam[property.ID] = param
				}
			}
		}
	}
	fmt.Printf("getParam=%+v\n",getParam)
	return nil
}

func (l *PublishLogic) HandleProperty(msg *types.Elements) error{
	l.Infof("PublishLogic|HandleProperty")
	dreq := dict.DeviceReq{}
	respTopic := fmt.Sprintf("$thing/down/property/%s/%s",l.topics[3],l.topics[4])
	err := utils.Unmarshal([]byte(msg.Payload), &dreq)
	if err != nil {
		return errors.Parameter.AddDetail("things topic is err:"+msg.Topic)
	}
	switch dreq.Method {
	case dict.REPORT:
		l.Infof("send topic=%s",respTopic)
		payload,_ := json.Marshal(dict.DeviceResp{
			Method: dict.REPORT_REPLY,
			ClientToken:dreq.ClientToken}.AddStatus(errors.OK))
		l.svcCtx.Mqtt.Publish(respTopic,0,false,payload)
	case dict.REPORT_INFO:
	case dict.GET_STATUS:
	default:
		return errors.Method
	}
	return nil
}

func (l *PublishLogic) HandleEvent(msg *types.Elements) error{
	l.Infof("PublishLogic|HandleEvent")
	return nil
}
func (l *PublishLogic) HandleAction(msg *types.Elements) error{
	l.Infof("PublishLogic|HandleAction")
	return nil
}

func (l *PublishLogic) HandleThing(msg *types.Elements) error{
	l.Infof("PublishLogic|HandleThing")
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
	l.Infof("PublishLogic|HandleOta")
	return nil
}

func (l *PublishLogic) HandleDefault(msg *types.Elements) error{
	l.Infof("PublishLogic|HandleDefault")
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
