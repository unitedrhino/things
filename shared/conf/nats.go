package conf

type NatsConf struct {
	Url   string `json:",default=nats://127.0.0.1:4222"` //nats的连接url
	User  string `json:",optional"`                      //用户名
	Pass  string `json:",optional"`                      //密码
	Token string `json:",optional"`
}
