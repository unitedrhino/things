package conf

type Email struct {
	From     string // 发件人  你自己要发邮件的邮箱
	Host     string // 服务器地址 例如 smtp.qq.com  请前往QQ或者你要发邮件的邮箱查看其smtp协议
	Secret   string // 密钥    用于登录的密钥 最好不要用邮箱密码 去邮箱smtp申请一个用于登录的密钥
	Nickname string // 昵称    发件人昵称 通常为自己的邮箱
	Port     int    `json:",default=465"` // 端口     请前往QQ或者你要发邮件的邮箱查看其smtp协议 大多为 465
	IsSSL    bool   // 是否SSL   是否开启SSL
}
