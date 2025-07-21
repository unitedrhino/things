package topics

const (
	DmDeviceInfoUnbind = "server.things.dm.device.info.unbind"
	DmDeviceInfoCreate = "server.things.dm.device.info.create"
	DmDeviceInfoDelete = "server.things.dm.device.info.delete"
	DmDeviceInfoUpdate = "server.things.dm.device.info.update"
	DmDeviceInfoBind   = "server.things.dm.device.info.bind"

	DmDeviceOnlineStatusChange = "server.things.dm.device.onlineStatus.change"
	DmDeviceStaticOneHour      = "server.things.dm.device.static.2Hour"     //2小时统计
	DmDeviceStaticHalfHour     = "server.things.dm.device.static.halfHour"  //半小时统计
	DmDeviceStaticOneMinute    = "server.things.dm.device.static.oneMinute" //1分钟统计
	DmProductInfoCreate        = "server.things.dm.product.info.create"
	DmProductInfoUpdate        = "server.things.dm.product.info.update"
	DmProductInfoDelete        = "server.things.dm.product.info.delete"
	DmProductCustomUpdate      = "server.things.dm.product.custom.update"   //产品脚本有更新
	DmOtaDeviceUpgradePush     = "server.things.dm.ota.device.upgrade.push" //ota设备推送
	DmOtaJobDelayRun           = "server.things.dm.ota.job.delay.run"       //任务延时启动
	// DmProtocolInfoUpdate 中间的是协议code
	DmProtocolInfoUpdate  = "server.things.dm.protocol.%s.update" //自定义协议配置有更新
	UdRuleTimer           = "server.things.ud.rule.timer"
	UdRuleTimerTenMinutes = "server.things.ud.rule.timer.tenMinutes"

	//最后一个参数是告警模式
	UdRuleAlarmNotify = "server.things.ud.rule.alarm.%s" //trigger:触发告警 relieve:解除告警
	DgOnlineTimer     = "server.things.dg.online.timer"
	DgOnlineTimer2    = "server.things.dg.online.timer2"

	PAliTimer = "server.things.pali.data.timer"

	DmActionCheckDelay = "server.things.dm.action.check.delay"
)

// 应用事件通知(设备状态变化,设备上报)
const (
	// ApplicationDeviceStatusConnected 设备登录状态推送 中间两个是产品id和设备名称
	ApplicationDeviceStatusConnected = "application.device.%s.%s.status.connected"
	// ApplicationDeviceStatusDisConnected 设备登出状态推送 中间两个是产品id和设备名称
	ApplicationDeviceStatusDisConnected = "application.device.%s.%s.status.disconnected"
	// ApplicationDeviceReportThingProperty 设备物模型属性上报通知 中间两个是产品id和设备名称,最后一个是属性id
	ApplicationDeviceReportThingProperty   = "application.device.%s.%s.report.thing.property.%s"
	ApplicationDeviceReportThingPropertyV2 = "application.v2.device.%s.%s.report.thing.property"
	// ApplicationDeviceReportThingEvent 设备物模型事件上报通知 中间两个是产品id和设备名称,最后两个是事件类型和事件id
	ApplicationDeviceReportThingEvent = "application.device.%s.%s.report.thing.event.%s.%s"
	// ApplicationDeviceReportThingAction 设备物模型事件上报通知 中间两个是产品id和设备名称,最后三个是actionID,请求类型(req resp)和调用方向
	ApplicationDeviceReportThingAction = "application.device.%s.%s.report.thing.action.%s.%s.%s"
	// ApplicationDeviceReportThingPropertyDevice 设备物模型属性上报通知 中间两个是产品id和设备名称
	ApplicationDeviceReportThingPropertyDevice = "application.device.%s.%s.report.thing.property"

	ApplicationDeviceReportThingEventAllDevice    = "application.device.*.*.report.thing.event.>"
	ApplicationDeviceReportThingPropertyAllDevice = "application.device.*.*.report.thing.property.>"
	ApplicationDeviceStatusConnectedAllDevice     = "application.device.*.*.status.connected"
	ApplicationDeviceStatusDisConnectedAllDevice  = "application.device.*.*.status.disconnected"
	ApplicationDeviceStatusAllDevice              = "application.device.*.*.status.>"
)
