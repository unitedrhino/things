package logic

import (
	"gitee.com/i-Things/share/def"
	"github.com/i-Things/things/service/viewsvr/internal/types"
)

func ToPageInfo(info *types.PageInfo) *def.PageInfo {
	if info == nil {
		return nil
	}
	return &def.PageInfo{
		Page: info.Page,
		Size: info.Size,
	}
}
