package userSubscribe

type UserSubscribe = string

var (
	DevicePropertyReport  = "devicePropertyReport"   //设备上报订阅
	DevicePropertyReport2 = "devicePropertyReportV2" //设备上报订阅
	DevicePublish         = "devicePublish"          //设备发布消息
	DeviceActionReport    = "deviceActionReport"     //设备行为消息
	DeviceEventReport     = "devicePropertyReport"   //设备上报订阅
	DeviceConn            = "deviceConn"             //设备连接消息
	DeviceOtaReport       = "deviceOtaReport"        //设备ota消息推送
)
