package topics

// dmsvr 发布的消息
const ()

// dmsvr发布的事件通知
const (
	DmDeviceInfoDelete         = "server.dm.device.info.delete"
	DmProductInfoDelete        = "server.dm.product.info.delete"
	DmProductCustomUpdate      = "server.dm.product.custom.update"       //产品脚本有更新
	DmProductSchemaUpdate      = "server.dm.product.schema.update"       //物模型有更新
	DmDeviceLogLevelUpdate     = "server.dm.device.logLevel.update"      //设备日志级别有更新
	DmDeviceGatewayUpdate      = "server.dm.device.gateway.update"       //网关下的子设备有改动
	DmDeviceRemoteConfigUpdate = "server.dm.product.remoteConfig.update" //远程配置推送
)
const (
	RuleSceneInfoUpdate = "server.rule.scene.info.update" //场景联动有修改
	RuleSceneInfoDelete = "server.rule.scene.info.delete" //场景联动有修改
)
