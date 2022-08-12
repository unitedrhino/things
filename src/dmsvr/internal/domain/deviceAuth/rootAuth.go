package deviceAuth

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/config"
)

func IsAdmin(white config.AuthWhite, ip string) bool {
	for _, whiteIp := range white.IpRange {
		if utils.MatchIP(ip, whiteIp) {
			return true
		}
	}
	return false
}
