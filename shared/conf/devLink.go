package conf

type DevLinkConf struct {
	Mode    string    `json:",default=mqtt"` //模式 默认mqtt
	SubMode string    `json:",default=emq"`  //
	Mqtt    *MqttConf `json:",optional"`
}
