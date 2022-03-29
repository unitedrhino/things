package third

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"time"
)

type DevClientConf struct {
	Mqtt MqttConf
}

type MqttConf struct {
	ClientID string   //在mqtt中的clientID
	Brokers  []string //mqtt服务器节点
	User     string   //用户名
	Pass     string   `json:",optional"` //密码
}

type Info struct {
	Timeout  time.Time
	SubTopic string
	ClientID string `json:"clientid"`
	Msg      chan *deviceTemplate.DeviceResp
}

func NewInfo(timeout time.Time, ClientID, SubTopic string) *Info {
	return &Info{
		Timeout:  timeout,
		ClientID: ClientID,
		Msg:      make(chan *deviceTemplate.DeviceResp, 1),
		SubTopic: SubTopic,
	}
}

func (i *Info) IsTimeOut() bool {
	return time.Now().Before(i.Timeout)
}

type DevClient struct {
	DeviceChan *utils.ExpMap
	client     mqtt.Client
}

func NewDevClient(conf DevClientConf) *DevClient {
	opts := mqtt.NewClientOptions()
	for _, broker := range conf.Mqtt.Brokers {
		opts.AddBroker(broker)
	}
	opts.SetClientID(conf.Mqtt.ClientID).SetUsername(conf.Mqtt.User).
		SetPassword(conf.Mqtt.Pass).SetAutoReconnect(true).SetConnectRetry(true)
	opts.OnConnect = func(client mqtt.Client) {
		//logx.Info("Connected")
	}
	mc := mqtt.NewClient(opts)
	mc.Connect()
	return &DevClient{
		client:     mc,
		DeviceChan: utils.NewExpMap(5 * time.Minute),
	}
}

func (d *DevClient) DeviceResp(Method, ClientToken string, topics []string, err error, data map[string]interface{}) {
	respMethod := deviceTemplate.GetMethod(Method)
	respTopic := fmt.Sprintf("%s/down/%s/%s/%s",
		topics[0], topics[2], topics[3], topics[4])
	payload, _ := json.Marshal(deviceTemplate.DeviceResp{
		Method:      respMethod,
		ClientToken: ClientToken,
		Data:        data}.AddStatus(err))
	d.client.Publish(respTopic, 0, false, payload)
}

func (d *DevClient) DeviceReq(ctx context.Context, req deviceTemplate.DeviceReq, pubTopic, subTopic string) (*deviceTemplate.DeviceResp, error) {
	payload, _ := json.Marshal(req)
	d.client.Publish(pubTopic, 1, false, payload)
	respInfo := NewInfo(time.Now().Add(5*time.Second), req.ClientToken, subTopic)
	d.DeviceChan.Map.Store(req.ClientToken, respInfo)
	defer d.DeviceChan.Map.Delete(req.ClientToken)
	for {
		select {
		case resp := <-respInfo.Msg:
			return resp, nil
		case <-ctx.Done():
			return &deviceTemplate.DeviceResp{}, errors.DeviceTimeOut
		}
	}
}

//服务器已经向设备发送了请求,这个是把回复送给对应的处理函数
func (d *DevClient) DeviceReqSendResp(resp *deviceTemplate.DeviceResp, topic string) error {
	c, ok := d.DeviceChan.Map.Load(resp.ClientToken)
	if ok != true {
		//如果没有找到,说明不是这个分区处理的或者这个消息超时了,不管是超时还是不是这个分区处理,由于无法判断是否是超时,所以会将返回固定错误码,让kafka转发给其他服务去处理,这里在很多服务的时候会出现高延迟,高性能损耗,以后需要进行优化
		return errors.Server
	}
	if c.(*Info).SubTopic != topic {
		return errors.RespParam.AddDetailf("device send topic:%s|sys topic:%s|not same", topic,
			c.(*Info).SubTopic)
	}
	c.(*Info).Msg <- resp
	return nil
}
