package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/device"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-uuid"
	"time"

	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type SendActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendActionLogic {
	return &SendActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}



func (l *SendActionLogic) SendAction(in *dm.SendActionReq) (*dm.SendActionResp, error) {
	SubTopic := fmt.Sprintf("$thing/up/action/%s/%s",in.ProductID,in.DeviceName)
	devMsg := make(chan []byte)
	token := l.svcCtx.Mqtt.Subscribe(SubTopic,0,func(c mqtt.Client,m mqtt.Message){
		devMsg<-m.Payload()
	})
	if !token.WaitTimeout(5*time.Second){
		l.Error("SendAction|Subscribe mqtt timeout")
		return nil, errors.System.AddDetail("SendAction|Subscribe mqtt timeout")
	}
	param := map[string]interface{}{}
	err := json.Unmarshal([]byte(in.InputParams),&param)
	if err != nil {
		return nil, errors.Parameter.AddDetail("SendAction|InputParams not right:",in.InputParams)
	}
	PubTopic := fmt.Sprintf("$thing/down/action/%s/%s",in.ProductID,in.DeviceName)
	uuid,err := uuid.GenerateUUID()
	if err != nil{
		l.Errorf("SendAction|GenerateUUID err:%v",err)
		return nil, errors.System.AddDetail(err)
	}
	payload, _ := json.Marshal(device.DeviceReq{
		Method:      "action",
		ClientToken: uuid,
		Timestamp: time.Now().Unix(),
		Params: param,})
	l.svcCtx.Mqtt.Publish(PubTopic,1,false,payload)

	return &dm.SendActionResp{}, nil
}
