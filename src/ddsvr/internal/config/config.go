package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	DevLink   DevLinkConf   //和设备交互的设置
	InnerLink InnerLinkConf //和things内部交互的设置
}

type DevLinkConf struct {
	Mode    string    `json:",default=mqtt"` //模式 默认mqtt
	SubMode string    `json:",default=emq"`  //
	Mqtt    *MqttConf `json:",optional"`
}
type MqttConf struct {
	ClientID string   //在mqtt中的clientID
	Brokers  []string //mqtt服务器节点
	User     string   //用户名
	Pass     string   `json:",optional"` //密码
}

type InnerLinkConf struct {
	Nats NatsConf
}
type NatsConf struct {
	Url   string `json:",default=nats://127.0.0.1:4222"` //nats的连接url
	User  string `json:",optional"`                      //用户名
	Pass  string `json:",optional"`                      //密码
	Token string `json:",optional"`
}
