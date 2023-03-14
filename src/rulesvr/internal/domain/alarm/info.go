package alarm

//告警记录状态
const (
	DealStateNone     = 1 //无告警
	DealStateAlarming = 2 //告警中
	DealStateAlarmed  = 3 //已处理
)

//告警配置级别
const (
	LevelInfo     = 1 //提醒
	LevelNormal   = 1 //一般
	LevelBad      = 1 //严重
	LevelUrgent   = 1 //紧急
	LevelSpUrgent = 1 //超级紧急
)
