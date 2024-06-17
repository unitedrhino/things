package logic

import (
	"gitee.com/i-Things/share/def"
	"github.com/i-Things/things/service/udsvr/pb/ud"
)

func ToPageInfo(info *ud.PageInfo) *def.PageInfo {
	if info == nil {
		return nil
	}
	return &def.PageInfo{
		Page: info.GetPage(),
		Size: info.GetSize(),
	}
}

func ToTimeRange(timeRange *ud.TimeRange) def.TimeRange {
	if timeRange == nil {
		return def.TimeRange{}
	}
	return def.TimeRange{Start: timeRange.Start, End: timeRange.End}
}
