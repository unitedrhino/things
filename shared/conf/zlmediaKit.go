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

type Gbsip struct {
	Lid            string `json:",optional"`
	Region         string `json:",optional"`
	Did            string `json:",optional"`
	Cid            string `json:",optional"`
	Dnum           int64  `json:",optional"`
	Cnum           int64  `json:",optional"`
	UDP            int64  `json:",optional"`
	DefaultDevName string `json:",optional"` // 国标设备可以直接配置   默认用户名：
	DefaultDevPswd string `json:",optional"` //                    默认密码
}
