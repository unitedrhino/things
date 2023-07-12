package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/atomic"
	"net/url"
	"os"
	"sync"
	"time"
)

const (
	broker      = "tcp://127.0.0.1:1883"
	pushPayload = `
{
  "method": "report",
  "clientToken": "%s",
  "params": {
    "bool": true,
    "int64": %v,
    "float64": 12.34234,
    "string": "%s",
    "position": {
      "longitude": 0.012334992177784441,
      "latitude": 50.045040130615234
    }
  }
}
`
	pushTopic = "$thing/up/property/26GOsI5N0vS/test%v"
	subTopic  = "$thing/down/property/26GOsI5N0vS/test%v"
	userName  = "ddsvr"
	password  = "iThings"
)

var (
	deviceTotal int64 = 1000
	pubInterval       = time.Second
	pubCount          = atomic.NewInt64(0)
	subCount          = atomic.NewInt64(0)
	errCount          = atomic.NewInt64(0)
	msgRwMutex        = sync.RWMutex{}
	msgMap            = map[string]time.Time{}
	delay       float64
)

func main() {
	args := os.Args
	if len(args) > 1 {
		deviceTotal = cast.ToInt64(args[1])
	}
	if len(args) > 2 {
		pubInterval = time.Duration(cast.ToInt(args[2])) * time.Second
	}
	fmt.Println(deviceTotal, pubInterval)
	var mcs = make([]mqtt.Client, 0, deviceTotal)
	for len(mcs) < int(deviceTotal) {
		mc, err := MqttInit(len(mcs))
		if err != nil {
			continue
		}
		mcs = append(mcs, mc)
	}
	logx.Infof("连接初始化完成")
	for i := int64(1); i <= deviceTotal; i++ {
		id := i
		mc := mcs[i-1]
		go Sub(id, mc)
		go Pub(id, mc)
	}
	var (
		pub int64
		sub int64
	)

	for true {
		time.Sleep(time.Second * 5)
		var unSubMsgCount int
		var before = time.Now
		func() {
			msgRwMutex.RLock()
			defer msgRwMutex.RUnlock()
			unSubMsgCount = len(msgMap)
		}()
		pubRate := pubCount.Load() - pub
		subRate := subCount.Load() - sub
		pub = pubCount.Load()
		sub = subCount.Load()
		useTime := time.Now().Sub(before())
		logx.WithDuration(useTime).Infof("errCount:%v pubCount:%v subCount:%v unSubMsgCount:%v delay:%.2fms pubRate:%v  subRate:%v",
			errCount.Load(), pubCount.Load(), subCount.Load(), unSubMsgCount, delay, float64(pubRate)/5, float64(subRate)/5)
	}
}

type CommonMsg struct { //消息内容通用字段
	Method      string `json:"method"`              //操作方法
	ClientToken string `json:"clientToken"`         //方便排查随机数
	Timestamp   int64  `json:"timestamp,omitempty"` //毫秒时间戳
	Code        int64  `json:"code,omitempty"`      //状态码
	Status      string `json:"status,omitempty"`    //返回信息
	Data        any    `json:"data,omitempty"`      //返回具体设备上报的最新数据内容
}

func Sub(id int64, mc mqtt.Client) {
	var err error
	topic := fmt.Sprintf(subTopic, id)
	err = mc.Subscribe(topic, 1, func(client mqtt.Client, message mqtt.Message) {
		var resp CommonMsg
		var ok bool
		var t time.Time
		var now = time.Now()
		json.Unmarshal(message.Payload(), &resp)
		func() {
			msgRwMutex.Lock()
			defer msgRwMutex.Unlock()
			t, ok = msgMap[resp.ClientToken]
			if ok {
				delete(msgMap, resp.ClientToken)
				delay = delay*0.8 + float64(now.Sub(t).Milliseconds())*0.2
			}
		}()
		if !ok { //已经获取过了或不是这里发的
			return
		}
		if resp.Code != errors.OK.Code {
			logx.Errorf("subscribe err id:%v topic:%v code:%v msg:%v", id, message.Topic(), resp.Code, resp.Status)
			errCount.Add(1)
		}
		subCount.Add(1)
	}).Error()
	logx.Must(err)
	for true {
		time.Sleep(time.Hour)
	}
}

func Pub(id int64, mc mqtt.Client) {
	defer func() {
		utils.Recover(context.Background())
		logx.Errorf("pub.id:%v finish", id)
	}()
	time.Sleep(time.Second * 2) //先让订阅订阅好再发布
	t := time.NewTicker(pubInterval)
	defer t.Stop()
	topic := fmt.Sprintf(pushTopic, id)
	i := 0
	for range t.C {
		i++
		tk, _ := uuid.GenerateUUID()
		tk = fmt.Sprintf("%s:%d", tk, i)
		payload := fmt.Sprintf(pushPayload, tk, i, time.Now().Format("2006-01-02 15:04:05.000"))
		err := mc.Publish(topic, 1, false, payload).WaitTimeout(time.Second * 5)
		if err != true {
			logx.Errorf("publish error err:%v id:%v", err, id)
		}
		pubCount.Add(1)
		func() {
			msgRwMutex.Lock()
			defer msgRwMutex.Unlock()
			msgMap[tk] = time.Now()
		}()
	}
}
func MqttInit(id int) (mc mqtt.Client, err error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	uuid, _ := uuid.GenerateUUID()
	clientID := fmt.Sprintf("pressureTest%v:%v", id, uuid)
	opts.SetClientID(clientID).SetUsername(userName).
		SetPassword(password).SetAutoReconnect(true).SetConnectRetry(true)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Infof("mqtt client:%v Connected ", id)
	})
	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		//logx.Infof("mqtt client connect attempt broker:%v", utils.Fmt(broker))
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt client connection lost err:%v", utils.Fmt(err))
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false {
		//logx.Info("mqtt Connect failure")
		err = fmt.Errorf("mqtt client connect failure id:%v", id)
		return
	}
	return mc, nil
}
