package conf

import (
	"github.com/spf13/cast"
)

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

func (counter WrongPasswordCounter) ParseWrongPassConf(userID string, ip string) []*LoginSafeCtlInfo {
	var res []*LoginSafeCtlInfo
	res = append(res, &LoginSafeCtlInfo{
		Prefix:  "login:wrongPassword:captcha:",
		Key:     "login:wrongPassword:captcha:" + userID,
		Timeout: 24 * 3600,
		Times:   counter.Captcha,
	})

	for i, v := range counter.Account {
		res = append(res, &LoginSafeCtlInfo{
			Prefix:    "login:wrongPassword:account:",
			Key:       "login:wrongPassword:account:" + cast.ToString(i+1) + ":" + userID,
			Timeout:   v.Statistics * 60,
			Times:     v.TriggerTimes,
			Forbidden: v.ForbiddenTime * 60,
		})
	}
	for i, v := range counter.Ip {
		res = append(res, &LoginSafeCtlInfo{
			Prefix:    "login:wrongPassword:ip:",
			Key:       "login:wrongPassword:ip:" + cast.ToString(i+1) + ":" + ip,
			Timeout:   v.Statistics * 60,
			Times:     v.TriggerTimes,
			Forbidden: v.ForbiddenTime * 60,
		})
	}

	return res
}
