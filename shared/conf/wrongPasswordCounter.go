package conf

type WrongPasswordCounter struct {
	Captcha int `json:",default=5"`
	Account []struct {
		Statistics    int `json:",default=1440"`
		TriggerTimes  int `json:",default=10"`
		ForbiddenTime int `json:",default=10"`
	}
	Ip []struct {
		Statistics    int `json:",default=1440"`
		TriggerTimes  int `json:",default=200"`
		ForbiddenTime int `json:",default=60"`
	}
}

type LoginSafeCtlInfo struct {
	Prefix    string // key前缀
	Key       string // redis key
	Timeout   int    // redis key 超时时间
	Times     int    // 错误密码次数
	Forbidden int    // 账号或ip冻结时间
}
