package logic

import (
	"gitee.com/i-Things/share/def"
	"github.com/i-Things/things/service/rulesvr/pb/rule"
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
