package conf

type MediaConf struct {
	Host   string `json:",optional"`
	Port   int64  `json:",optional"`
	Secret string `json:",optional"`
}
