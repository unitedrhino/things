package def

type LogLevel = int64

const (
	LogClose LogLevel = 1 //关闭
	LogError LogLevel = 2 //错误
	LogWarn  LogLevel = 3 //告警
	LogInfo  LogLevel = 4 //信息
	LogDebug LogLevel = 5 //调试
)

var LogLevelTextToIntMap = map[string]LogLevel{
	"关闭": LogClose,
	"错误": LogError,
	"告警": LogWarn,
	"信息": LogInfo,
	"调试": LogDebug,
}

const (
	GatewayBind   = 1 //绑定
	GatewayUnbind = 2 //解绑
)

type DeviceStatus = int64

const (
	DeviceStatusInactive DeviceStatus = 0 // 未激活
	DeviceStatusOnline   DeviceStatus = 1 //在线
	DeviceStatusOffline  DeviceStatus = 2 //离线
)

// 视频设备状态
const (
	StatusInactive int64 = 0 // 未激活
	StatusOnline   int64 = 1 //在线
	StatusOffline  int64 = 2 //离线
)
