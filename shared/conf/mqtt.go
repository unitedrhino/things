package conf

type MqttConf struct {
	ClientID string   //在mqtt中的clientID
	Brokers  []string //mqtt服务器节点
	User     string   `json:",default=root"` //用户名
	Pass     string   `json:",optional"`     //密码
}
