package eventBus

// 服务自己的消息
const (
	DmDeviceInfoDelete    = "server.dm.device.info.delete"
	DmProductInfoDelete   = "server.dm.product.info.delete"
	DmProductCustomUpdate = "server.dm.product.custom.update" //产品脚本有更新
	DmProductSchemaUpdate = "server.dm.product.schema.update" //物模型有更新
)
