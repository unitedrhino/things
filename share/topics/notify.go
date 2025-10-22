package topics

const (

	//下面的是纯服务自己用的
	DmOtaJobDelayRun           = "server.things.dm.ota.job.delay.run" //任务延时启动
	UdRuleTimer                = "server.things.ud.rule.timer"
	UdRuleTimerTenMinutes      = "server.things.ud.rule.timer.tenMinutes"
	DgOnlineTimer              = "server.things.dg.online.timer"
	DgOnlineTimer2             = "server.things.dg.online.timer2"
	DmActionCheckDelay         = "server.things.dm.action.check.delay"
	DmDeviceOnlineStatusChange = "server.things.dm.device.onlineStatus.change"
	DmDeviceStaticOneHour      = "server.things.dm.device.static.2Hour"     //2小时统计
	DmDeviceStaticHalfHour     = "server.things.dm.device.static.halfHour"  //半小时统计
	DmDeviceStaticOneMinute    = "server.things.dm.device.static.oneMinute" //1分钟统计
	DmOtaDeviceUpgradePush     = "server.things.dm.ota.device.upgrade.push" //ota设备推送
	DmProtocolInfoUpdate       = "server.things.dm.protocol.%s.update"      //自定义协议配置有更新 中间的是协议code

)

// 应用事件通知(设备状态变化,设备上报)
const (
	DmDeviceInfoUnbind  = "app.%s.things.dm.device.info.unbind"
	DmDeviceInfoCreate  = "app.%s.things.dm.device.info.create"
	DmDeviceInfoDelete  = "app.%s.things.dm.device.info.delete"
	DmDeviceInfoUpdate  = "app.%s.things.dm.device.info.update"
	DmDeviceInfoBind    = "app.%s.things.dm.device.info.bind"
	DmProductInfoCreate = "app.%s.things.dm.product.info.create"
	DmProductInfoUpdate = "app.%s.things.dm.product.info.update"
	DmProductInfoDelete = "app.%s.things.dm.product.info.delete"
	//DmProductCustomUpdate  = "app.%s.things.dm.product.custom.update"   //产品脚本有更新

	UdRuleAlarmNotify = "app.%s.things.ud.rule.alarm.%s" //最后一个参数是告警模式 trigger:触发告警 relieve:解除告警

	// ApplicationDeviceStatusConnected 设备登录状态推送 中间两个是产品id和设备名称
	ApplicationDeviceStatusConnected = "app.%s.device.%s.%s.status.connected"
	// ApplicationDeviceStatusDisConnected 设备登出状态推送 中间两个是产品id和设备名称
	ApplicationDeviceStatusDisConnected = "app.%s.device.%s.%s.status.disconnected"
	// ApplicationDeviceReportThingProperty 设备物模型属性上报通知 中间两个是产品id和设备名称,最后一个是属性id
	ApplicationDeviceReportThingProperty   = "app.%s.device.%s.%s.report.thing.property.%s"
	ApplicationDeviceReportThingPropertyV2 = "app.%s.v2.device.%s.%s.report.thing.property"
	// ApplicationDeviceReportThingEvent 设备物模型事件上报通知 中间两个是产品id和设备名称,最后两个是事件类型和事件id
	ApplicationDeviceReportThingEvent = "app.%s.device.%s.%s.report.thing.event.%s.%s"
	// ApplicationDeviceReportThingAction 设备物模型行为上报通知 中间两个是产品id和设备名称,最后三个是actionID,请求类型(req resp)和调用方向
	ApplicationDeviceReportThingAction = "app.%s.device.%s.%s.report.thing.action.%s.%s.%s"
	// ApplicationDeviceReportThingPropertyDevice 设备物模型属性上报通知 中间两个是产品id和设备名称
	ApplicationDeviceReportThingPropertyDevice = "app.%s.device.%s.%s.report.thing.property"

	ApplicationDeviceReportThingEventAllDevice    = "app.%s.device.*.*.report.thing.event.>"
	ApplicationDeviceReportThingPropertyAllDevice = "app.%s.device.*.*.report.thing.property.>"
	ApplicationDeviceStatusConnectedAllDevice     = "app.%s.device.*.*.status.connected"
	ApplicationDeviceStatusDisConnectedAllDevice  = "app.%s.device.*.*.status.disconnected"
	ApplicationDeviceStatusAllDevice              = "app.%s.device.*.*.status.>"
)
