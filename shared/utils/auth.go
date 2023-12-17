package utils

import "github.com/i-Things/things/shared/conf"

// Auth 在名单内返回true
func Auth(a conf.AuthConf, userName, password, ipaddr string) bool {
	var userCompare bool
	for _, user := range a.Users {
		if userName == user.UserName {
			userCompare = false
			if password == user.Password {
				userCompare = true
			}
			break
		}
	}
	if !userCompare {
		return false
	}
	if len(a.IpRange) == 0 {
		//如果没有,表示不开启ip白名单模式
		return true
	}
	for _, whiteIp := range a.IpRange {
		if MatchIP(ipaddr, whiteIp) {
			return true
		}
	}
	return false
}
