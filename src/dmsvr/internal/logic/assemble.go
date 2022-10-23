package logic

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
)

func ToPageInfo(info *dm.PageInfo) def.PageInfo {
	return def.PageInfo{
		Page: info.GetPage(),
		Size: info.GetSize(),
	}
}
