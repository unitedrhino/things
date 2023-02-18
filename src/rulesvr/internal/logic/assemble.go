package logic

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/rulesvr/pb/rule"
)

func ToPageInfo(info *rule.PageInfo) *def.PageInfo {
	if info == nil {
		return nil
	}
	return &def.PageInfo{
		Page: info.GetPage(),
		Size: info.GetSize(),
	}
}
