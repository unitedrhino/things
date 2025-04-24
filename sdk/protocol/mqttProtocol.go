package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/events/topics"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/clients"
	"gitee.com/unitedrhino/things/share/devices"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/timex"
	"math"

	"strings"
	"time"
)

type DevHandle func(ctx context.Context, topic string, payload []byte) error

type ConnHandle func(ctx context.Context, info devices.DevConn) (devices.DevConn, error)

type MqttProtocol struct {
	*CoreProtocol
	MqttClient   *clients.MqttClient
	DevSubHandle map[string]DevHandle
	ConnHandle   ConnHandle
}

func NewMqttProtocol(c conf.EventConf, pi *dm.ProtocolInfo, pc *CoreProtocolConf, mqttc conf.DevLinkConf) (*MqttProtocol, error) {
	if mqttc.Mqtt == nil {
		return nil, errors.Parameter.AddMsg("DevLinkConf need")
	}
	mc, err := clients.NewMqttClient(mqttc.Mqtt)
	if err != nil {
		return nil, err
	}
	lightProto, err := NewCoreProtocol(c, pi, pc)
	if err != nil {
		return nil, err
	}
	return &MqttProtocol{CoreProtocol: lightProto, MqttClient: mc, DevSubHandle: make(map[string]DevHandle)}, nil
}

//// RegisterRpc 如果不想自己维护proto,则实现方法注入即可
//func (m *MqttProtocol) RegisterRpc(c zrpc.RpcServerConf, handle dg.DeviceAuthServer) {
//	m.RegisterRpcServer(c, func(grpcServer *grpc.Server) {
//		dg.RegisterDeviceAuthServer(grpcServer, handle)
//	})
//}

func (m *MqttProtocol) SubscribeDevMsg(topic string, handle DevHandle) error {
	newTopic := fmt.Sprintf("$share/%s/%s", m.ServerName, topic)
	m.DevSubHandle[newTopic] = handle
	return nil
}

func (m *MqttProtocol) PublishToDev(ctx context.Context, topic string, payload []byte) error {
	err := m.MqttClient.Publish(topic, 1, false, payload)
	if err != nil {
		logx.WithContext(ctx).Errorf("PublishToDev failure err:%v topic:%v", err, topic)
	} else {
		logx.WithContext(ctx).Infof("PublishToDev  topic:%v payload:%v", topic, string(payload))
	}
	return err
}

func (m *MqttProtocol) Start() error {
	err := m.subscribes(nil)
	if err != nil {
		return err
	}
	clients.SetMqttSetOnConnectHandler(func(cli mqtt.Client) {
		err := m.subscribes(cli)
		if err != nil {
			logx.Errorf("%s.mqttSetOnConnectHandler.subscribes err:%v", utils.FuncName(), err)
		}
	})
	err = m.CoreProtocol.Start()
	if err != nil {
		return err
	}
	return nil
}

func (m *MqttProtocol) subscribes(cli mqtt.Client) error {
	for topic, handle := range m.DevSubHandle {
		err := m.subscribeWithFunc(cli, topic, handle)
		if err != nil {
			return err
		}
	}
	return nil
}

type (
	//登录登出消息
	ConnectMsg struct {
		UserName string `json:"username"`
		Ts       int64  `json:"ts"`
		Address  string `json:"ipaddress"`
		ClientID string `json:"clientid"`
		Reason   string `json:"reason"`
	}
)

const (
	DeviceMqttDevice   = "device:mqtt:device"
	DeviceMqttClientID = "device:mqtt:clientID"
	DeviceLastActivity = "device:lastActivity"
)

func GenDeviceTopicKey(dev devices.Core) string {
	return fmt.Sprintf("%v:%v", dev.ProductID, dev.DeviceName)
}

func UpdateDeviceActivity(ctx context.Context, dev devices.Core) {
	_, err := caches.GetStore().Zadd(DeviceLastActivity, time.Now().Unix(), GenDeviceTopicKey(dev))
	if err != nil {
		logx.WithContext(ctx).Info(err)
	}
}

func UpdatesDeviceActivity(ctx context.Context, devs []devices.Core) {
	if len(devs) == 0 {
		return
	}
	var ps []redis.Pair
	now := time.Now().Unix()
	for _, v := range devs {
		ps = append(ps, redis.Pair{
			Key:   GenDeviceTopicKey(v),
			Score: now,
		})
	}
	_, err := caches.GetStore().Zadds(DeviceLastActivity, ps...)
	if err != nil {
		logx.WithContext(ctx).Info(err)
	}
}

func GetActivityDevices(ctx context.Context) (map[devices.Core]struct{}, error) {
	devs, err := caches.GetStore().ZrangeCtx(ctx, DeviceLastActivity, math.MinInt64, math.MaxInt64)
	if err != nil {
		return nil, err
	}
	var ret = make(map[devices.Core]struct{}, len(devs))
	for _, v := range devs {
		ps := strings.Split(v, ":")
		ret[devices.Core{ProductID: ps[0], DeviceName: ps[1]}] = struct{}{}
	}
	return ret, nil
}

func DeleteDeviceActivity(ctx context.Context, dev devices.Core) {
	_, err := caches.GetStore().Zrem(DeviceLastActivity, GenDeviceTopicKey(dev))
	if err != nil {
		logx.WithContext(ctx).Info(err)
	}
}

func (m *MqttProtocol) SubscribeDevConn(handle ConnHandle) error {
	m.ConnHandle = handle
	newTopic := fmt.Sprintf("$share/%s/%s", m.ServerName, "$SYS/brokers/+/clients/#")
	m.DevSubHandle[newTopic] = func(ctx context.Context, topic string, payload []byte) error {
		var (
			msg ConnectMsg
			err error
		)
		err = json.Unmarshal(payload, &msg)
		if err != nil {
			logx.Error(err)
			return err
		}
		do := devices.DevConn{
			UserName:  msg.UserName,
			Timestamp: msg.Ts, //毫秒时间戳
			Address:   msg.Address,
			ClientID:  msg.ClientID,
			Reason:    msg.Reason,
		}
		if strings.HasSuffix(topic, "/disconnected") {
			logx.WithContext(ctx).Infof("%s.disconnected topic:%v message:%v err:%v",
				utils.FuncName(), topic, string(payload), err)
			do.Action = devices.ActionDisconnected
		} else {
			do.Action = devices.ActionConnected
			logx.WithContext(ctx).Infof("%s.connected topic:%v message:%v err:%v",
				utils.FuncName(), topic, string(payload), err)
		}
		newDo, err := handle(ctx, do)
		if err != nil {
			return nil //不是该类型的设备
		}
		// 忽略下线后会导致设备正常上下线也会被忽略,这里先不处理
		//if do.Reason == "takeovered" || do.Reason == "takenover" || do.Reason == "discard" || do.Reason == "discarded" { //连接还在的时候被别人顶了,忽略这种下线
		//	logx.Errorf("忽略的下线状态:%v", utils.Fmt(do))
		//	return nil
		//}
		dev := devices.Core{
			ProductID:  newDo.ProductID,
			DeviceName: newDo.DeviceName,
		}
		err = caches.GetStore().Hset(DeviceMqttDevice, GenDeviceTopicKey(dev), utils.MarshalNoErr(do))
		if err != nil {
			logx.Error(err)
		}
		err = caches.GetStore().Hset(DeviceMqttClientID, do.ClientID, utils.MarshalNoErr(dev))
		if err != nil {
			logx.Error(err)
		}
		newDo.ClientID = fmt.Sprintf("%s&%s", newDo.ProductID, newDo.DeviceName)
		switch do.Action {
		case devices.ActionConnected:
			err = m.FastEvent.Publish(ctx, topics.DeviceUpStatusConnected, newDo)
		case devices.ActionDisconnected:
			err = m.FastEvent.Publish(ctx, topics.DeviceUpStatusDisconnected, newDo)
		default:
			panic("not support conn type")
		}
		if err != nil {
			logx.Errorf("%s.publish  err:%v", utils.FuncName(), err)
			return err
		}
		return err
	}
	return nil
}

func (m *MqttProtocol) subscribeWithFunc(cli mqtt.Client, topic string, handle func(ctx context.Context, topic string, payload []byte) error) error {
	return m.MqttClient.Subscribe(cli, topic,
		1, func(client mqtt.Client, message mqtt.Message) {
			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
			ctx, span := ctxs.StartSpan(ctx, message.Topic(), "")
			utils.Go(ctx, func() {
				defer cancel()
				//dgsvr 订阅到了设备端数据，此时调用StartSpan方法，将订阅到的主题推送给jaeger
				//此时的ctx已经包含当前节点的span信息，会随着 handle(ctx).Publish 传递到下个节点
				defer span.End()
				startTime := timex.Now()
				duration := timex.Since(startTime)
				ctx = ctxs.WithRoot(ctx)
				err := handle(ctx, message.Topic(), message.Payload())
				if err != nil {
					logx.WithContext(ctx).Errorf("%s.handle failure err:%v topic:%v payload:%v", utils.FuncName(), err, topic, string(message.Payload()))
				}
				logx.WithContext(ctx).WithDuration(duration).Infof(
					"SubscribeDevicePublish topic:%v message:%v err:%v",
					message.Topic(), string(message.Payload()), err)
			})
		})
}
