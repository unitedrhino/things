package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/google/uuid"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cast"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	mqttInitOnce sync.Once
	mqttClient   *MqttClient
	// mqttSetOnConnectHandler 如果会话断开可以通过该回调函数来重新订阅消息
	//不使用mqtt的clean session是因为会话保持期间共享订阅也会给离线的客户端,这会导致在线的客户端丢失消息
	mqttSetOnConnectHandler func(cli mqtt.Client)
)

type MqttClient struct {
	clients []mqtt.Client
	cfg     *conf.MqttConf
}

func NewMqttClient(conf *conf.MqttConf) (mcs *MqttClient, err error) {
	mqttInitOnce.Do(func() {
		var clients []mqtt.Client
		for len(clients) < conf.ConnNum {
			var (
				mc mqtt.Client
			)
			var tryTime = 5
			for i := tryTime; i > 0; i-- {
				mc, err = initMqtt(conf)
				if err != nil { //出现并发情况的时候可能联犀的http还没启动完毕
					logx.Errorf("mqtt 连接失败 重试剩余次数:%v", i-1)
					time.Sleep(time.Second * time.Duration(tryTime) / time.Duration(i))
					continue
				}
				break
			}
			if err != nil {
				logx.Errorf("mqtt 连接失败 conf:%#v  err:%v", conf, err)
				os.Exit(-1)
			}
			clients = append(clients, mc)
			var cli = MqttClient{clients: clients, cfg: conf}
			mqttClient = &cli

		}
	})
	return mqttClient, err
}

func SetMqttSetOnConnectHandler(f func(cli mqtt.Client)) {
	mqttSetOnConnectHandler = f
}

type EmqResp struct {
	Code    string `json:"code"` //如果不在线返回: CLIENTID_NOT_FOUND
	Message string `json:"message"`
}

type MutSubReq struct {
	Topic string `json:"topic"`
	Qos   int    `json:"qos"`
	Nl    int    `json:"nl"`
	Rap   int    `json:"rap"`
	Rh    int    `json:"rh"`
}

func (m MqttClient) SetClientMutSub(ctx context.Context, clientID string, topics []string, qos int) error {
	logx.WithContext(ctx).Infof("SetClientMut clientID:%v,topics:%v", clientID, topics)
	ts, err := m.GetClientSub(ctx, clientID)
	if err != nil {
		return err
	}
	var topicSet = map[string]struct{}{}
	for _, topic := range topics {
		topicSet[topic] = struct{}{}
	}
	for _, t := range ts {
		if _, ok := topicSet[t]; ok {
			delete(topicSet, t)
		}
	}
	topics = utils.SetToSlice(topicSet)
	if len(topics) == 0 {
		return nil
	}
	if m.cfg.OpenApi == nil {
		return errors.System.AddMsg("未开启登录检查")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	var ret []*MutSubReq
	var req []*MutSubReq
	for _, v := range topics {
		req = append(req, &MutSubReq{
			Topic: v,
			Qos:   qos,
		})
	}
	var errs []error
	var body []byte
	var tryTime = 5
	for i := tryTime; i > 0; i-- {
		_, body, errs = greq.Post(fmt.Sprintf("%s/api/v5/clients/%s/subscribe/bulk", oa.Host,
			url.QueryEscape(clientID))).Send(&req).EndStruct(&ret)
		if errs != nil {
			time.Sleep(time.Second / 2 * time.Duration(tryTime) / time.Duration(i))
			continue
		}
		break
	}
	if errs != nil {
		return errors.System.AddDetail(errs, string(body))
	}

	return nil
}
func (m MqttClient) GetClientSub(ctx context.Context, clientID string) ([]string, error) {
	logx.WithContext(ctx).Infof("GetClientSub clientID:%v", clientID)
	if m.cfg.OpenApi == nil {
		return nil, errors.System.AddMsg("未开启登录检查")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	var ret []*MutSubReq
	var errs []error
	var body []byte
	var tryTime = 5
	for i := tryTime; i > 0; i-- {
		_, body, errs = greq.Get(fmt.Sprintf("%s/api/v5/clients/%s/subscriptions", oa.Host,
			url.QueryEscape(clientID))).EndStruct(&ret)
		if errs != nil {
			time.Sleep(time.Second / 2 * time.Duration(tryTime) / time.Duration(i))
			continue
		}
		break
	}
	if errs != nil {
		return nil, errors.System.AddDetail(errs, string(body))
	}
	var topics []string
	for _, v := range ret {
		topics = append(topics, v.Topic)
	}
	return topics, nil
}

func (m MqttClient) SetClientMutUnSub(ctx context.Context, clientID string, topics []string) error {
	logx.WithContext(ctx).Infof("SetClientMut clientID:%v,topics:%v", clientID, topics)
	if m.cfg.OpenApi == nil {
		return errors.System.AddMsg("未开启登录检查")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	var ret []*MutSubReq
	var req []*MutSubReq
	for _, v := range topics {
		req = append(req, &MutSubReq{
			Topic: v,
			Qos:   0,
		})
	}
	var errs []error
	var body []byte
	var tryTime = 5
	for i := tryTime; i > 0; i-- {
		_, body, errs = greq.Post(fmt.Sprintf("%s/api/v5/clients/%s/unsubscribe/bulk", oa.Host,
			url.QueryEscape(clientID))).Send(&req).EndStruct(&ret)
		if errs != nil {
			time.Sleep(time.Second / 2 * time.Duration(tryTime) / time.Duration(i))
			continue
		}
		break
	}
	if errs != nil {
		return errors.System.AddDetail(errs, string(body))
	}
	return nil
}

// https://www.emqx.io/docs/zh/v5.5/admin/api-docs.html#tag/Clients/paths/~1clients~1%7Bclientid%7D/get
func (m MqttClient) CheckIsOnline(ctx context.Context, clientID string) (bool, error) {
	if m.cfg.OpenApi == nil {
		return false, errors.System.AddMsg("未开启登录检查")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	var ret EmqResp
	resp, rets, errs := greq.Get(fmt.Sprintf("%s/api/v5/clients/%s", oa.Host, url.QueryEscape(clientID))).EndStruct(&ret)
	if errs != nil {
		return false, errors.System.AddDetail(errs)
	}
	if resp.StatusCode != http.StatusOK {
		return false, errors.System.AddDetail(string(rets))
	}
	if errs != nil {
		return false, errors.System.AddDetail(errs)
	}
	if ret.Code == "" {
		return true, nil
	}
	return false, nil
}

type EmqGetClientsResp struct {
	Data []struct {
		HeapSize                         int       `json:"heap_size"`
		SendMsgDroppedExpired            int       `json:"send_msg.dropped.expired"`
		SendOct                          int       `json:"send_oct"`
		RecvMsgQos1                      int       `json:"recv_msg.qos1"`
		IsPersistent                     bool      `json:"is_persistent"`
		SendPkt                          int       `json:"send_pkt"`
		CleanStart                       bool      `json:"clean_start"`
		InflightCnt                      int       `json:"inflight_cnt"`
		Node                             string    `json:"node"`
		SendMsgDroppedQueueFull          int       `json:"send_msg.dropped.queue_full"`
		AwaitingRelCnt                   int       `json:"awaiting_rel_cnt"`
		InflightMax                      int       `json:"inflight_max"`
		CreatedAt                        time.Time `json:"created_at"`
		SubscriptionsCnt                 int       `json:"subscriptions_cnt"`
		MailboxLen                       int       `json:"mailbox_len"`
		SendCnt                          int       `json:"send_cnt"`
		Connected                        bool      `json:"connected"`
		IpAddress                        string    `json:"ip_address"`
		AwaitingRelMax                   int       `json:"awaiting_rel_max"`
		RecvMsgQos2                      int       `json:"recv_msg.qos2"`
		ProtoVer                         int       `json:"proto_ver"`
		Mountpoint                       string    `json:"mountpoint"`
		ProtoName                        string    `json:"proto_name"`
		Port                             int       `json:"port"`
		ConnectedAt                      time.Time `json:"connected_at"`
		EnableAuthn                      bool      `json:"enable_authn"`
		ExpiryInterval                   int       `json:"expiry_interval"`
		Username                         *string   `json:"username"`
		RecvMsg                          int       `json:"recv_msg"`
		RecvOct                          int       `json:"recv_oct"`
		SendMsgDroppedTooLarge           int       `json:"send_msg.dropped.too_large"`
		Keepalive                        int       `json:"keepalive"`
		SendMsgQos1                      int       `json:"send_msg.qos1"`
		SendMsgQos2                      int       `json:"send_msg.qos2"`
		RecvMsgQos0                      int       `json:"recv_msg.qos0"`
		SendMsgQos0                      int       `json:"send_msg.qos0"`
		SubscriptionsMax                 string    `json:"subscriptions_max"`
		MqueueMax                        int       `json:"mqueue_max"`
		MqueueDropped                    int       `json:"mqueue_dropped"`
		Clientid                         string    `json:"clientid"`
		IsBridge                         bool      `json:"is_bridge"`
		Peerport                         int       `json:"peerport"`
		SendMsg                          int       `json:"send_msg"`
		Listener                         string    `json:"listener"`
		RecvCnt                          int       `json:"recv_cnt"`
		RecvPkt                          int       `json:"recv_pkt"`
		RecvMsgDropped                   int       `json:"recv_msg.dropped"`
		SendMsgDropped                   int       `json:"send_msg.dropped"`
		RecvMsgDroppedAwaitPubrelTimeout int       `json:"recv_msg.dropped.await_pubrel_timeout"`
		Reductions                       int       `json:"reductions"`
		MqueueLen                        int       `json:"mqueue_len"`
	} `json:"data"`
	Meta struct {
		Count   int64 `json:"count"`
		Hasnext bool  `json:"hasnext"`
		Limit   int   `json:"limit"`
		Page    int   `json:"page"`
	} `json:"meta"`
}

type OnlineClientsInfo struct {
	ClientID    string
	UserName    string
	ConnectedAt time.Time
}
type GetOnlineClientsFilter struct {
	UserName string
}

type PageInfo struct {
	Page int64 `json:"page" form:"page"`         // 页码
	Size int64 `json:"pageSize" form:"pageSize"` // 每页大小
}

func (m MqttClient) GetOnlineClients(ctx context.Context, f GetOnlineClientsFilter, page *PageInfo) ([]*devices.DevConn, int64, error) {
	if m.cfg.OpenApi == nil {
		return nil, 0, errors.System.AddMsg("未开启登录检查")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	greq.Get(fmt.Sprintf("%s/api/v5/clients", oa.Host))
	if f.UserName != "" {
		greq.Query("username=" + f.UserName)
	}
	if page != nil {
		greq.Query(fmt.Sprintf("page=%v", page.Page))
		greq.Query(fmt.Sprintf("limit=%v", page.Size))
	}
	var ret EmqGetClientsResp
	resp, rets, errs := greq.EndStruct(&ret)
	if errs != nil {
		return nil, 0, errors.System.AddDetail(errs)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, errors.System.AddDetail(string(rets))
	}
	var infos []*devices.DevConn
	for _, v := range ret.Data {
		infos = append(infos, &devices.DevConn{
			ClientID:  v.Clientid,
			UserName:  cast.ToString(v.Username),
			Timestamp: v.ConnectedAt.UnixMilli(),
		})
	}
	return infos, ret.Meta.Count, nil
}

func (m MqttClient) Subscribe(cli mqtt.Client, topic string, qos byte, callback mqtt.MessageHandler) error {
	var clients = m.clients
	if cli != nil {
		clients = []mqtt.Client{cli}
	}
	for _, c := range clients {
		err := c.Subscribe(topic, qos, callback).Error()
		if err != nil {
			return errors.System.AddDetail(err)
		}
	}
	return nil
}

func (m MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	id := rand.Intn(len(m.clients))
	return m.clients[id].Publish(topic, qos, retained, payload).Error()
}

func initMqtt(conf *conf.MqttConf) (mc mqtt.Client, err error) {
	opts := mqtt.NewClientOptions()
	for _, broker := range conf.Brokers {
		opts.AddBroker(broker)
	}
	uuid := uuid.NewString()
	opts.SetClientID(conf.ClientID + "/" + uuid).SetUsername(conf.User).SetPassword(conf.Pass)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logx.Info("mqtt client Connected")
		if mqttSetOnConnectHandler != nil {
			mqttSetOnConnectHandler(client)
		}
	})

	opts.SetAutoReconnect(true).SetMaxReconnectInterval(30 * time.Second) //意外离线的重连参数
	opts.SetConnectRetry(true).SetConnectRetryInterval(5 * time.Second)   //首次连接的重连参数

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logx.Infof("mqtt 正在尝试连接 broker:%v", utils.Fmt(broker))
		return tlsCfg
	})
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logx.Errorf("mqtt 连接丢失 err:%v", utils.Fmt(err))
	})
	mc = mqtt.NewClient(opts)
	er2 := mc.Connect().WaitTimeout(5 * time.Second)
	if er2 == false {
		logx.Info("mqtt 连接失败")
		err = fmt.Errorf("mqtt 连接失败")
		return
	}
	return
}
