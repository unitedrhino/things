package conf

type MediaConf struct {
	Host   string `json:",optional"`
	Port   int64  `json:",optional"`
	Secret string `json:",optional"`
}

type External struct {
	Host string `json:",default=0.0.0.0,optional"`
	Port int64  `json:",optional"`
}
