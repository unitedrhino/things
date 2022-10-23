package deviceAuth

import (
	"github.com/i-Things/things/shared/utils"
)

type (
	AuthWhite struct {
		Users   []string //用户白名单
		IpRange []string //ip 及ip段
	}
)

func IsAdmin(white AuthWhite, info AuthInfo) bool {
	var userCompare bool
	for _, user := range white.Users {
		if info.Username == user {
			userCompare = true
			break
		}
	}
	if !userCompare {
		return false
	}
	for _, whiteIp := range white.IpRange {
		if utils.MatchIP(info.Ip, whiteIp) {
			return true
		}
	}
	return false
}
