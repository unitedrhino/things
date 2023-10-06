package conf

const (
	EventModeNats   = "nats"   //使用nats来通讯
	EventModeNatsJs = "natsJs" //使用nats的jetstream来通讯
	EventModeDirect = "direct" //直接调用
)

type EventConf struct {
	Mode string `json:",default=nats,options=nats|natsJs|direct"`
	Nats NatsConf
}
