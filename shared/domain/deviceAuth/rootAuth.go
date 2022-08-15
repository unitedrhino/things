package deviceAuth

import (
	"github.com/i-Things/things/shared/utils"
)

type AuthWhite struct {
	IpRange []string //ip 及ip段
}

func IsAdmin(white AuthWhite, ip string) bool {
	for _, whiteIp := range white.IpRange {
		if utils.MatchIP(ip, whiteIp) {
			return true
		}
	}
	return false
}
