package device

type LOG_LEVEL = int64

const (
	LOG_CLOSE LOG_LEVEL = 1 //关闭
	LOG_ERROR LOG_LEVEL = 2 //错误
	LOG_WARN  LOG_LEVEL = 3 //告警
	LOG_INFO  LOG_LEVEL = 4 //信息
	LOG_DEBUG LOG_LEVEL = 5 //调试
)
