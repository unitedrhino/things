package logic

import (
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/udsvr/pb/ud"
)

func ToPageInfo(info *ud.PageInfo) *stores.PageInfo {
	return utils.Copy[stores.PageInfo](info)
}

func ToTimeRange(timeRange *ud.TimeRange) *def.TimeRange {
	if timeRange == nil {
		return nil
	}
	return &def.TimeRange{Start: timeRange.Start, End: timeRange.End}
}
