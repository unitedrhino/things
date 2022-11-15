package def

type LogLevel = int64

const (
	LogClose LogLevel = 1 //关闭
	LogError LogLevel = 2 //错误
	LogWarn  LogLevel = 3 //告警
	LogInfo  LogLevel = 4 //信息
	LogDebug LogLevel = 5 //调试
)

const (
	GatewayBind   = 1 //绑定
	GatewayUnbind = 2 //解绑
)
