package domain

import (
	"gitee.com/unitedrhino/share/utils"
	"time"
)

func GenLastRunTime(t time.Time, execAt int64) time.Time {
	dt := utils.DaySecToTime(t, execAt)
	if dt.Before(t) { //如果执行时间在之前,则今天不需要执行了
		return utils.GetEndTime(t)
	}
	//如果还没有执行,则修改为当天的最早的时间,到了执行时间就会自动执行
	return utils.GetZeroTime(t)
}
