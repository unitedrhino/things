package conf

const (
	InnerLinkModeNats   = "nats"   //使用nats来通讯
	InnerLinkModeDirect = "direct" //直接调用
)

type InnerLinkConf struct {
	Mode string `json:",default=nats,options=nats|direct"`
	Nats NatsConf
}
