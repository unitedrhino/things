package productInfo

type Auth = int64

const (
	AuthPwd  Auth = 1 //账密认证
	AuthCert Auth = 2 //证书认证
)

type Net = int64

const (
	NetOther Net = 1 //其他
	NetWifi  Net = 2 //wi-fi
	NetG234  Net = 3 //2G/3G/4G
	NetG5    Net = 4 //5G
	NetBle   Net = 5 //蓝牙
	NetLora  Net = 6 //LoRaWAN
)

type Dev = int64

const (
	DevDevice  Dev = 1 //设备
	DevGateway Dev = 2 //网关
	DevSubset  Dev = 3 //子设备
)

type Data = int64

const (
	DataUnknown  Data = 0 //禁止为此参数
	DataCustom   Data = 1 //自定义
	DataTemplate Data = 2 //数据模板
)

type AutoReg = int64

const (
	AutoRegClose AutoReg = 1 //关闭
	AutoRegOpen  AutoReg = 2 //打开
	AutoRegAuto  AutoReg = 3 //打开并自动创建设备
)
