package conf

type AuthConf struct {
	IpRange []string `json:",optional"` //白名单ip 及ip段
	Users   []AuthUserInfo
}

type AuthUserInfo struct {
	UserName string // 内部为服务名
	Password string // 密码
}
