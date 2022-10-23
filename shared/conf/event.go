package conf

const (
	EventModeNats   = "nats"   //使用nats来通讯
	EventModeDirect = "direct" //直接调用
)

type EventConf struct {
	Mode string `json:",default=nats,options=nats|direct"`
	Nats NatsConf
}
