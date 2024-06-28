package logic

import (
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/stores"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToPageInfo(info *ud.PageInfo) *stores.PageInfo {
	if info == nil {
		return nil
	}
	return &stores.PageInfo{
		Page: info.GetPage(),
		Size: info.GetSize(),
	}
}

func ToTimeRange(timeRange *ud.TimeRange) *def.TimeRange {
	if timeRange == nil {
		return nil
	}
	return &def.TimeRange{Start: timeRange.Start, End: timeRange.End}
}
