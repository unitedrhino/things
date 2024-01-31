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
	Dnum           int32  `json:",optional"`
	Cnum           int32  `json:",optional"`
	NetT           string `json:",optional"` //使用UDP或者TCP
	Host           string `json:",optional"` //gb28181平台IP
	Port           int32  `json:",optional"` //gb28181使用的UDP端口
	DefaultDevName string `json:",optional"` // 国标设备可以直接配置   默认用户名：
	DefaultDevPswd string `json:",optional"` //                    默认密码
}
