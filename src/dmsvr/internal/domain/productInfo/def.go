package productInfo

type AUTH = int64

const (
	AuthPwd  AUTH = 1 //账密认证
	AuthCert AUTH = 2 //证书认证
)

type NET = int64

const (
	NetOther NET = 1 //其他
	NetWifi  NET = 2 //wi-fi
	NetG234  NET = 3 //2G/3G/4G
	NetG5    NET = 4 //5G
	NetBle   NET = 5 //蓝牙
	NetLora  NET = 6 //LoRaWAN
)

type DEV = int64

const (
	DevDevice  DEV = 1 //设备
	DevGateway DEV = 2 //网关
	DevSubset  DEV = 3 //子设备
)

type DATA = int64

const (
	DataUnknown  DATA = 0 //禁止为此参数
	DataCustom   DATA = 1 //自定义
	DataTemplate DATA = 2 //数据模板
)

type AutoReg = int64

const (
	AutoRegClose AutoReg = 1 //关闭
	AutoRegOpen  AutoReg = 2 //打开
	AutoRegAuto  AutoReg = 3 //打开并自动创建设备
)
