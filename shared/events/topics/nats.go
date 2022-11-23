package topics

// 设备交互相关topic
const (

	// DeviceUpThing 物模型 最后两个是产品id和设备名称
	DeviceUpThing    = "device.up.thing.%s.%s"
	DeviceUpThingAll = "device.up.thing.>"

	// DeviceUpGateway 网关与子设备 最后两个是产品id和设备名称
	DeviceUpGateway    = "device.up.gateway.%s.%s"
	DeviceUpGatewayAll = "device.up.gateway.>"

	// DeviceUpOta ota升级相关 最后两个是产品id和设备名称
	DeviceUpOta    = "device.up.ota.%s.%s"
	DeviceUpOtaAll = "device.up.ota.>"

	// DeviceUpShadow 设备影子  最后两个是产品id和设备名称
	DeviceUpShadow    = "device.up.shadow.%s.%s"
	DeviceUpShadowAll = "device.up.shadow.>"

	// DeviceUpConfig 设备远程配置 最后两个是产品id和设备名称
	DeviceUpConfig    = "device.up.config.%s.%s"
	DeviceUpConfigAll = "device.up.config.>"

	// DeviceUpSDKLog 设备调试日志 最后两个是产品id和设备名称
	DeviceUpSDKLog    = "device.up.sdkLog.%s.%s"
	DeviceUpSDKLogAll = "device.up.sdkLog.>"

	// DeviceUpStatusConnected 设备登录后向内部推送以下topic
	DeviceUpStatusConnected = "device.up.status.connected"
	// DeviceUpStatusDisconnected 设备的登出后向内部推送以下topic
	DeviceUpStatusDisconnected = "device.up.status.disconnected"

	// DeviceDownAll dd模块订阅以下topic,收到内部的发布消息后向设备推送
	DeviceDownAll = "device.down"
)

//dmsvr发布的事件通知
const (
	DmProductUpdateSchema      = "server.dm.product.update.schema"       //物模型有更新
	DmDeviceUpdateLogLevel     = "server.dm.device.update.logLevel"      //设备日志级别有更新
	DmDeviceUpdateGateway      = "server.dm.device.update.gateway"       //网关下的子设备有改动
	DmDeviceUpdateRemoteConfig = "server.dm.product.update.remoteConfig" //远程配置推送
)
