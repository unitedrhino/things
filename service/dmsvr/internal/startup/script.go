package startup

import (
	"gitee.com/unitedrhino/things/service/dmsvr/client/deviceinteract"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	deviceinteractServer "gitee.com/unitedrhino/things/service/dmsvr/internal/server/deviceinteract"
	devicemanageServer "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"reflect"
)

func ScriptSymbolInit(svcCtx *svc.ServiceContext) {
	svcCtx.ScriptTrans.AddSymbol("dm/dm", dmSymbolInit(svcCtx))
	return
}

func dmSymbolInit(svcCtx *svc.ServiceContext) map[string]reflect.Value {
	AbnormalLogIndexReq := dm.AbnormalLogIndexReq{}
	AbnormalLogIndexResp := dm.AbnormalLogIndexResp{}
	AbnormalLogInfo := dm.AbnormalLogInfo{}
	ActionRespReq := dm.ActionRespReq{}
	ActionSendReq := dm.ActionSendReq{}
	ActionSendResp := dm.ActionSendResp{}
	CommonSchemaCreateReq := dm.CommonSchemaCreateReq{}
	CommonSchemaIndexReq := dm.CommonSchemaIndexReq{}
	CommonSchemaIndexResp := dm.CommonSchemaIndexResp{}
	CommonSchemaInfo := dm.CommonSchemaInfo{}
	CommonSchemaUpdateReq := dm.CommonSchemaUpdateReq{}
	CompareInt64 := dm.CompareInt64{}
	CompareString := dm.CompareString{}
	CustomTopic := dm.CustomTopic{}
	DeviceBindTokenInfo := dm.DeviceBindTokenInfo{}
	DeviceBindTokenReadReq := dm.DeviceBindTokenReadReq{}
	DeviceCore := dm.DeviceCore{}
	DeviceCountInfo := dm.DeviceCountInfo{}
	DeviceCountReq := dm.DeviceCountReq{}
	DeviceCountResp := dm.DeviceCountResp{}
	DeviceError := dm.DeviceError{}
	DeviceGatewayBindDevice := dm.DeviceGatewayBindDevice{}
	DeviceGatewayIndexReq := dm.DeviceGatewayIndexReq{}
	DeviceGatewayIndexResp := dm.DeviceGatewayIndexResp{}
	DeviceGatewayMultiCreateReq := dm.DeviceGatewayMultiCreateReq{}
	DeviceGatewayMultiSaveReq := dm.DeviceGatewayMultiSaveReq{}
	DeviceGatewaySign := dm.DeviceGatewaySign{}
	DeviceGroupMultiSaveReq := dm.DeviceGroupMultiSaveReq{}
	DeviceInfo := dm.DeviceInfo{}
	DeviceInfoBindReq := dm.DeviceInfoBindReq{}
	DeviceInfoCanBindReq := dm.DeviceInfoCanBindReq{}
	DeviceInfoCount := dm.DeviceInfoCount{}
	DeviceInfoCountReq := dm.DeviceInfoCountReq{}
	DeviceInfoDeleteReq := dm.DeviceInfoDeleteReq{}
	DeviceInfoIndexReq := dm.DeviceInfoIndexReq{}
	DeviceInfoIndexResp := dm.DeviceInfoIndexResp{}
	DeviceInfoMultiBindReq := dm.DeviceInfoMultiBindReq{}
	DeviceInfoMultiBindResp := dm.DeviceInfoMultiBindResp{}
	DeviceInfoMultiUpdateReq := dm.DeviceInfoMultiUpdateReq{}
	DeviceInfoReadReq := dm.DeviceInfoReadReq{}
	DeviceInfoUnbindReq := dm.DeviceInfoUnbindReq{}
	DeviceModuleVersion := dm.DeviceModuleVersion{}
	DeviceModuleVersionIndexReq := dm.DeviceModuleVersionIndexReq{}
	DeviceModuleVersionIndexResp := dm.DeviceModuleVersionIndexResp{}
	DeviceModuleVersionReadReq := dm.DeviceModuleVersionReadReq{}
	DeviceMoveReq := dm.DeviceMoveReq{}
	DeviceOnlineMultiFix := dm.DeviceOnlineMultiFix{}
	DeviceOnlineMultiFixReq := dm.DeviceOnlineMultiFixReq{}
	DeviceProfile := dm.DeviceProfile{}
	DeviceProfileIndexReq := dm.DeviceProfileIndexReq{}
	DeviceProfileIndexResp := dm.DeviceProfileIndexResp{}
	DeviceProfileReadReq := dm.DeviceProfileReadReq{}
	DeviceResetReq := dm.DeviceResetReq{}
	DeviceSchema := dm.DeviceSchema{}
	DeviceSchemaIndexReq := dm.DeviceSchemaIndexReq{}
	DeviceSchemaIndexResp := dm.DeviceSchemaIndexResp{}
	DeviceSchemaMultiCreateReq := dm.DeviceSchemaMultiCreateReq{}
	DeviceSchemaMultiDeleteReq := dm.DeviceSchemaMultiDeleteReq{}
	DeviceSchemaTslReadReq := dm.DeviceSchemaTslReadReq{}
	DeviceSchemaTslReadResp := dm.DeviceSchemaTslReadResp{}
	DeviceShareInfo := dm.DeviceShareInfo{}
	DeviceTransferReq := dm.DeviceTransferReq{}
	DeviceTypeCountReq := dm.DeviceTypeCountReq{}
	DeviceTypeCountResp := dm.DeviceTypeCountResp{}
	EdgeSendReq := dm.EdgeSendReq{}
	EdgeSendResp := dm.EdgeSendResp{}
	Empty := dm.Empty{}
	EventLogIndexReq := dm.EventLogIndexReq{}
	EventLogIndexResp := dm.EventLogIndexResp{}
	EventLogInfo := dm.EventLogInfo{}
	Firmware := dm.Firmware{}
	FirmwareFile := dm.FirmwareFile{}
	FirmwareInfo := dm.FirmwareInfo{}
	FirmwareInfoDeleteReq := dm.FirmwareInfoDeleteReq{}
	FirmwareInfoDeleteResp := dm.FirmwareInfoDeleteResp{}
	FirmwareInfoIndexReq := dm.FirmwareInfoIndexReq{}
	FirmwareInfoIndexResp := dm.FirmwareInfoIndexResp{}
	FirmwareInfoReadReq := dm.FirmwareInfoReadReq{}
	FirmwareInfoReadResp := dm.FirmwareInfoReadResp{}
	FirmwareResp := dm.FirmwareResp{}
	GatewayCanBindIndexReq := dm.GatewayCanBindIndexReq{}
	GatewayCanBindIndexResp := dm.GatewayCanBindIndexResp{}
	GatewayGetFoundReq := dm.GatewayGetFoundReq{}
	GatewayNotifyBindSendReq := dm.GatewayNotifyBindSendReq{}
	GroupCore := dm.GroupCore{}
	GroupDeviceMultiDeleteReq := dm.GroupDeviceMultiDeleteReq{}
	GroupDeviceMultiSaveReq := dm.GroupDeviceMultiSaveReq{}
	GroupInfo := dm.GroupInfo{}
	GroupInfoCreateReq := dm.GroupInfoCreateReq{}
	GroupInfoIndexReq := dm.GroupInfoIndexReq{}
	GroupInfoIndexResp := dm.GroupInfoIndexResp{}
	GroupInfoMultiCreateReq := dm.GroupInfoMultiCreateReq{}
	GroupInfoReadReq := dm.GroupInfoReadReq{}
	GroupInfoUpdateReq := dm.GroupInfoUpdateReq{}
	HubLogIndexReq := dm.HubLogIndexReq{}
	HubLogIndexResp := dm.HubLogIndexResp{}
	HubLogInfo := dm.HubLogInfo{}
	IDPath := dm.IDPath{}
	IDPathWithUpdate := dm.IDPathWithUpdate{}
	OtaFirmwareDeviceCancelReq := dm.OtaFirmwareDeviceCancelReq{}
	OtaFirmwareDeviceConfirmReq := dm.OtaFirmwareDeviceConfirmReq{}
	OtaFirmwareDeviceIndexReq := dm.OtaFirmwareDeviceIndexReq{}
	OtaFirmwareDeviceIndexResp := dm.OtaFirmwareDeviceIndexResp{}
	OtaFirmwareDeviceInfo := dm.OtaFirmwareDeviceInfo{}
	OtaFirmwareDeviceRetryReq := dm.OtaFirmwareDeviceRetryReq{}
	OtaFirmwareFile := dm.OtaFirmwareFile{}
	OtaFirmwareFileIndexReq := dm.OtaFirmwareFileIndexReq{}
	OtaFirmwareFileIndexResp := dm.OtaFirmwareFileIndexResp{}
	OtaFirmwareFileInfo := dm.OtaFirmwareFileInfo{}
	OtaFirmwareFileReq := dm.OtaFirmwareFileReq{}
	OtaFirmwareFileResp := dm.OtaFirmwareFileResp{}
	OtaFirmwareInfo := dm.OtaFirmwareInfo{}
	OtaFirmwareInfoCreateReq := dm.OtaFirmwareInfoCreateReq{}
	OtaFirmwareInfoIndexReq := dm.OtaFirmwareInfoIndexReq{}
	OtaFirmwareInfoIndexResp := dm.OtaFirmwareInfoIndexResp{}
	OtaFirmwareInfoUpdateReq := dm.OtaFirmwareInfoUpdateReq{}
	OtaFirmwareJobIndexReq := dm.OtaFirmwareJobIndexReq{}
	OtaFirmwareJobIndexResp := dm.OtaFirmwareJobIndexResp{}
	OtaFirmwareJobInfo := dm.OtaFirmwareJobInfo{}
	OtaJobByDeviceIndexReq := dm.OtaJobByDeviceIndexReq{}
	OtaJobDynamicInfo := dm.OtaJobDynamicInfo{}
	OtaJobStaticInfo := dm.OtaJobStaticInfo{}
	OtaModuleInfo := dm.OtaModuleInfo{}
	OtaModuleInfoIndexReq := dm.OtaModuleInfoIndexReq{}
	OtaModuleInfoIndexResp := dm.OtaModuleInfoIndexResp{}
	PageInfo := dm.PageInfo{}
	PageInfo_OrderBy := dm.PageInfo_OrderBy{}
	Point := dm.Point{}
	ProductCategory := dm.ProductCategory{}
	ProductCategoryIndexReq := dm.ProductCategoryIndexReq{}
	ProductCategoryIndexResp := dm.ProductCategoryIndexResp{}
	ProductCategorySchemaIndexReq := dm.ProductCategorySchemaIndexReq{}
	ProductCategorySchemaIndexResp := dm.ProductCategorySchemaIndexResp{}
	ProductCategorySchemaMultiSaveReq := dm.ProductCategorySchemaMultiSaveReq{}
	ProductCustom := dm.ProductCustom{}
	ProductCustomReadReq := dm.ProductCustomReadReq{}
	ProductCustomUi := dm.ProductCustomUi{}
	ProductInfo := dm.ProductInfo{}
	ProductInfoDeleteReq := dm.ProductInfoDeleteReq{}
	ProductInfoIndexReq := dm.ProductInfoIndexReq{}
	ProductInfoIndexResp := dm.ProductInfoIndexResp{}
	ProductInfoReadReq := dm.ProductInfoReadReq{}
	ProductInitReq := dm.ProductInitReq{}
	ProductRemoteConfig := dm.ProductRemoteConfig{}
	ProductSchemaCreateReq := dm.ProductSchemaCreateReq{}
	ProductSchemaDeleteReq := dm.ProductSchemaDeleteReq{}
	ProductSchemaIndexReq := dm.ProductSchemaIndexReq{}
	ProductSchemaIndexResp := dm.ProductSchemaIndexResp{}
	ProductSchemaInfo := dm.ProductSchemaInfo{}
	ProductSchemaMultiCreateReq := dm.ProductSchemaMultiCreateReq{}
	ProductSchemaTslImportReq := dm.ProductSchemaTslImportReq{}
	ProductSchemaTslReadReq := dm.ProductSchemaTslReadReq{}
	ProductSchemaTslReadResp := dm.ProductSchemaTslReadResp{}
	ProductSchemaUpdateReq := dm.ProductSchemaUpdateReq{}
	PropertyControlMultiSendReq := dm.PropertyControlMultiSendReq{}
	PropertyControlMultiSendResp := dm.PropertyControlMultiSendResp{}
	PropertyControlSendMsg := dm.PropertyControlSendMsg{}
	PropertyControlSendReq := dm.PropertyControlSendReq{}
	PropertyControlSendResp := dm.PropertyControlSendResp{}
	PropertyGetReportMultiSendReq := dm.PropertyGetReportMultiSendReq{}
	PropertyGetReportMultiSendResp := dm.PropertyGetReportMultiSendResp{}
	PropertyGetReportSendMsg := dm.PropertyGetReportSendMsg{}
	PropertyGetReportSendReq := dm.PropertyGetReportSendReq{}
	PropertyGetReportSendResp := dm.PropertyGetReportSendResp{}
	PropertyLogIndexReq := dm.PropertyLogIndexReq{}
	PropertyLogIndexResp := dm.PropertyLogIndexResp{}
	PropertyLogInfo := dm.PropertyLogInfo{}
	PropertyLogLatestIndexReq := dm.PropertyLogLatestIndexReq{}
	ProtocolConfigField := dm.ProtocolConfigField{}
	ProtocolConfigInfo := dm.ProtocolConfigInfo{}
	ProtocolInfo := dm.ProtocolInfo{}
	ProtocolInfoIndexReq := dm.ProtocolInfoIndexReq{}
	ProtocolInfoIndexResp := dm.ProtocolInfoIndexResp{}
	ProtocolService := dm.ProtocolService{}
	ProtocolServiceIndexReq := dm.ProtocolServiceIndexReq{}
	ProtocolServiceIndexResp := dm.ProtocolServiceIndexResp{}
	RemoteConfigCreateReq := dm.RemoteConfigCreateReq{}
	RemoteConfigIndexReq := dm.RemoteConfigIndexReq{}
	RemoteConfigIndexResp := dm.RemoteConfigIndexResp{}
	RemoteConfigLastReadReq := dm.RemoteConfigLastReadReq{}
	RemoteConfigLastReadResp := dm.RemoteConfigLastReadResp{}
	RemoteConfigPushAllReq := dm.RemoteConfigPushAllReq{}
	RespReadReq := dm.RespReadReq{}
	RootCheckReq := dm.RootCheckReq{}
	SdkLogIndexReq := dm.SdkLogIndexReq{}
	SdkLogIndexResp := dm.SdkLogIndexResp{}
	SdkLogInfo := dm.SdkLogInfo{}
	SendLogIndexReq := dm.SendLogIndexReq{}
	SendLogIndexResp := dm.SendLogIndexResp{}
	SendLogInfo := dm.SendLogInfo{}
	SendMsgReq := dm.SendMsgReq{}
	SendMsgResp := dm.SendMsgResp{}
	SendOption := dm.SendOption{}
	ShadowIndex := dm.ShadowIndex{}
	ShadowIndexResp := dm.ShadowIndexResp{}
	SharePerm := dm.SharePerm{}
	StatusLogIndexReq := dm.StatusLogIndexReq{}
	StatusLogIndexResp := dm.StatusLogIndexResp{}
	StatusLogInfo := dm.StatusLogInfo{}
	TimeRange := dm.TimeRange{}
	UserDeviceCollectSave := dm.UserDeviceCollectSave{}
	UserDeviceShareIndexReq := dm.UserDeviceShareIndexReq{}
	UserDeviceShareIndexResp := dm.UserDeviceShareIndexResp{}
	UserDeviceShareInfo := dm.UserDeviceShareInfo{}
	UserDeviceShareMultiAcceptReq := dm.UserDeviceShareMultiAcceptReq{}
	UserDeviceShareMultiDeleteReq := dm.UserDeviceShareMultiDeleteReq{}
	UserDeviceShareMultiInfo := dm.UserDeviceShareMultiInfo{}
	UserDeviceShareMultiToken := dm.UserDeviceShareMultiToken{}
	UserDeviceShareReadReq := dm.UserDeviceShareReadReq{}
	WithID := dm.WithID{}
	WithIDChildren := dm.WithIDChildren{}
	WithIDCode := dm.WithIDCode{}
	WithProfile := dm.WithProfile{}
	PublishMsg := deviceMsg.PublishMsg{}
	return map[string]reflect.Value{
		"PublishMsg":                 reflect.ValueOf(PublishMsg),
		"ActionSend":                 reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).ActionSend),
		"ActionRead":                 reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).ActionRead),
		"ActionResp":                 reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).ActionResp),
		"PropertyGetReportSend":      reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).PropertyGetReportSend),
		"PropertyGetReportMultiSend": reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).PropertyGetReportMultiSend),
		"PropertyControlSend":        reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).PropertyControlSend),
		"PropertyControlMultiSend":   reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).PropertyControlMultiSend),
		"PropertyControlRead":        reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).PropertyControlRead),
		"GatewayGetFoundSend":        reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).GatewayGetFoundSend),
		"GatewayNotifyBindSend":      reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx)).GatewayNotifyBindSend),

		"DeviceInfoCreate":         reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoCreate),
		"DeviceInfoUpdate":         reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoUpdate),
		"DeviceInfoDelete":         reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoDelete),
		"DeviceInfoIndex":          reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoIndex),
		"DeviceInfoMultiUpdate":    reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoMultiUpdate),
		"DeviceInfoRead":           reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoRead),
		"DeviceInfoBind":           reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoBind),
		"DeviceBindTokenRead":      reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceBindTokenRead),
		"DeviceBindTokenCreate":    reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceBindTokenCreate),
		"DeviceInfoMultiBind":      reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoMultiBind),
		"DeviceInfoCanBind":        reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoCanBind),
		"DeviceInfoUnbind":         reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceInfoUnbind),
		"DeviceTransfer":           reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceTransfer),
		"DeviceReset":              reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceReset),
		"DeviceMove":               reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceMove),
		"DeviceModuleVersionRead":  reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceModuleVersionRead),
		"DeviceModuleVersionIndex": reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceModuleVersionIndex),
		"DeviceGatewayMultiCreate": reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceGatewayMultiCreate),
		"DeviceGatewayMultiUpdate": reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceGatewayMultiUpdate),
		"DeviceGatewayIndex":       reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceGatewayIndex),
		"DeviceGatewayMultiDelete": reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceGatewayMultiDelete),
		"DeviceProfileRead":        reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceProfileRead),
		"DeviceProfileDelete":      reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceProfileDelete),
		"DeviceProfileUpdate":      reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceProfileUpdate),
		"DeviceProfileIndex":       reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceProfileIndex),
		"DeviceSchemaUpdate":       reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceSchemaUpdate),
		"DeviceSchemaCreate":       reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceSchemaCreate),
		"DeviceSchemaMultiCreate":  reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceSchemaMultiCreate),
		"DeviceSchemaMultiDelete":  reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceSchemaMultiDelete),
		"DeviceSchemaIndex":        reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceSchemaIndex),
		"DeviceSchemaTslRead":      reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceSchemaTslRead),
		"DeviceGroupMultiCreate":   reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceGroupMultiCreate),
		"DeviceGroupMultiUpdate":   reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceGroupMultiUpdate),
		"DeviceGroupMultiDelete":   reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx)).DeviceGroupMultiDelete),

		"AbnormalLogIndexReq":               reflect.ValueOf(AbnormalLogIndexReq),
		"AbnormalLogIndexResp":              reflect.ValueOf(AbnormalLogIndexResp),
		"AbnormalLogInfo":                   reflect.ValueOf(AbnormalLogInfo),
		"ActionRespReq":                     reflect.ValueOf(ActionRespReq),
		"ActionSendReq":                     reflect.ValueOf(ActionSendReq),
		"ActionSendResp":                    reflect.ValueOf(ActionSendResp),
		"CommonSchemaCreateReq":             reflect.ValueOf(CommonSchemaCreateReq),
		"CommonSchemaIndexReq":              reflect.ValueOf(CommonSchemaIndexReq),
		"CommonSchemaIndexResp":             reflect.ValueOf(CommonSchemaIndexResp),
		"CommonSchemaInfo":                  reflect.ValueOf(CommonSchemaInfo),
		"CommonSchemaUpdateReq":             reflect.ValueOf(CommonSchemaUpdateReq),
		"CompareInt64":                      reflect.ValueOf(CompareInt64),
		"CompareString":                     reflect.ValueOf(CompareString),
		"CustomTopic":                       reflect.ValueOf(CustomTopic),
		"DeviceBindTokenInfo":               reflect.ValueOf(DeviceBindTokenInfo),
		"DeviceBindTokenReadReq":            reflect.ValueOf(DeviceBindTokenReadReq),
		"DeviceCore":                        reflect.ValueOf(DeviceCore),
		"DeviceCountInfo":                   reflect.ValueOf(DeviceCountInfo),
		"DeviceCountReq":                    reflect.ValueOf(DeviceCountReq),
		"DeviceCountResp":                   reflect.ValueOf(DeviceCountResp),
		"DeviceError":                       reflect.ValueOf(DeviceError),
		"DeviceGatewayBindDevice":           reflect.ValueOf(DeviceGatewayBindDevice),
		"DeviceGatewayIndexReq":             reflect.ValueOf(DeviceGatewayIndexReq),
		"DeviceGatewayIndexResp":            reflect.ValueOf(DeviceGatewayIndexResp),
		"DeviceGatewayMultiCreateReq":       reflect.ValueOf(DeviceGatewayMultiCreateReq),
		"DeviceGatewayMultiSaveReq":         reflect.ValueOf(DeviceGatewayMultiSaveReq),
		"DeviceGatewaySign":                 reflect.ValueOf(DeviceGatewaySign),
		"DeviceGroupMultiSaveReq":           reflect.ValueOf(DeviceGroupMultiSaveReq),
		"DeviceInfo":                        reflect.ValueOf(DeviceInfo),
		"DeviceInfoBindReq":                 reflect.ValueOf(DeviceInfoBindReq),
		"DeviceInfoCanBindReq":              reflect.ValueOf(DeviceInfoCanBindReq),
		"DeviceInfoCount":                   reflect.ValueOf(DeviceInfoCount),
		"DeviceInfoCountReq":                reflect.ValueOf(DeviceInfoCountReq),
		"DeviceInfoDeleteReq":               reflect.ValueOf(DeviceInfoDeleteReq),
		"DeviceInfoIndexReq":                reflect.ValueOf(DeviceInfoIndexReq),
		"DeviceInfoIndexResp":               reflect.ValueOf(DeviceInfoIndexResp),
		"DeviceInfoMultiBindReq":            reflect.ValueOf(DeviceInfoMultiBindReq),
		"DeviceInfoMultiBindResp":           reflect.ValueOf(DeviceInfoMultiBindResp),
		"DeviceInfoMultiUpdateReq":          reflect.ValueOf(DeviceInfoMultiUpdateReq),
		"DeviceInfoReadReq":                 reflect.ValueOf(DeviceInfoReadReq),
		"DeviceInfoUnbindReq":               reflect.ValueOf(DeviceInfoUnbindReq),
		"DeviceModuleVersion":               reflect.ValueOf(DeviceModuleVersion),
		"DeviceModuleVersionIndexReq":       reflect.ValueOf(DeviceModuleVersionIndexReq),
		"DeviceModuleVersionIndexResp":      reflect.ValueOf(DeviceModuleVersionIndexResp),
		"DeviceModuleVersionReadReq":        reflect.ValueOf(DeviceModuleVersionReadReq),
		"DeviceMoveReq":                     reflect.ValueOf(DeviceMoveReq),
		"DeviceOnlineMultiFix":              reflect.ValueOf(DeviceOnlineMultiFix),
		"DeviceOnlineMultiFixReq":           reflect.ValueOf(DeviceOnlineMultiFixReq),
		"DeviceProfile":                     reflect.ValueOf(DeviceProfile),
		"DeviceProfileIndexReq":             reflect.ValueOf(DeviceProfileIndexReq),
		"DeviceProfileIndexResp":            reflect.ValueOf(DeviceProfileIndexResp),
		"DeviceProfileReadReq":              reflect.ValueOf(DeviceProfileReadReq),
		"DeviceResetReq":                    reflect.ValueOf(DeviceResetReq),
		"DeviceSchema":                      reflect.ValueOf(DeviceSchema),
		"DeviceSchemaIndexReq":              reflect.ValueOf(DeviceSchemaIndexReq),
		"DeviceSchemaIndexResp":             reflect.ValueOf(DeviceSchemaIndexResp),
		"DeviceSchemaMultiCreateReq":        reflect.ValueOf(DeviceSchemaMultiCreateReq),
		"DeviceSchemaMultiDeleteReq":        reflect.ValueOf(DeviceSchemaMultiDeleteReq),
		"DeviceSchemaTslReadReq":            reflect.ValueOf(DeviceSchemaTslReadReq),
		"DeviceSchemaTslReadResp":           reflect.ValueOf(DeviceSchemaTslReadResp),
		"DeviceShareInfo":                   reflect.ValueOf(DeviceShareInfo),
		"DeviceTransferReq":                 reflect.ValueOf(DeviceTransferReq),
		"DeviceTypeCountReq":                reflect.ValueOf(DeviceTypeCountReq),
		"DeviceTypeCountResp":               reflect.ValueOf(DeviceTypeCountResp),
		"EdgeSendReq":                       reflect.ValueOf(EdgeSendReq),
		"EdgeSendResp":                      reflect.ValueOf(EdgeSendResp),
		"Empty":                             reflect.ValueOf(Empty),
		"EventLogIndexReq":                  reflect.ValueOf(EventLogIndexReq),
		"EventLogIndexResp":                 reflect.ValueOf(EventLogIndexResp),
		"EventLogInfo":                      reflect.ValueOf(EventLogInfo),
		"Firmware":                          reflect.ValueOf(Firmware),
		"FirmwareFile":                      reflect.ValueOf(FirmwareFile),
		"FirmwareInfo":                      reflect.ValueOf(FirmwareInfo),
		"FirmwareInfoDeleteReq":             reflect.ValueOf(FirmwareInfoDeleteReq),
		"FirmwareInfoDeleteResp":            reflect.ValueOf(FirmwareInfoDeleteResp),
		"FirmwareInfoIndexReq":              reflect.ValueOf(FirmwareInfoIndexReq),
		"FirmwareInfoIndexResp":             reflect.ValueOf(FirmwareInfoIndexResp),
		"FirmwareInfoReadReq":               reflect.ValueOf(FirmwareInfoReadReq),
		"FirmwareInfoReadResp":              reflect.ValueOf(FirmwareInfoReadResp),
		"FirmwareResp":                      reflect.ValueOf(FirmwareResp),
		"GatewayCanBindIndexReq":            reflect.ValueOf(GatewayCanBindIndexReq),
		"GatewayCanBindIndexResp":           reflect.ValueOf(GatewayCanBindIndexResp),
		"GatewayGetFoundReq":                reflect.ValueOf(GatewayGetFoundReq),
		"GatewayNotifyBindSendReq":          reflect.ValueOf(GatewayNotifyBindSendReq),
		"GroupCore":                         reflect.ValueOf(GroupCore),
		"GroupDeviceMultiDeleteReq":         reflect.ValueOf(GroupDeviceMultiDeleteReq),
		"GroupDeviceMultiSaveReq":           reflect.ValueOf(GroupDeviceMultiSaveReq),
		"GroupInfo":                         reflect.ValueOf(GroupInfo),
		"GroupInfoCreateReq":                reflect.ValueOf(GroupInfoCreateReq),
		"GroupInfoIndexReq":                 reflect.ValueOf(GroupInfoIndexReq),
		"GroupInfoIndexResp":                reflect.ValueOf(GroupInfoIndexResp),
		"GroupInfoMultiCreateReq":           reflect.ValueOf(GroupInfoMultiCreateReq),
		"GroupInfoReadReq":                  reflect.ValueOf(GroupInfoReadReq),
		"GroupInfoUpdateReq":                reflect.ValueOf(GroupInfoUpdateReq),
		"HubLogIndexReq":                    reflect.ValueOf(HubLogIndexReq),
		"HubLogIndexResp":                   reflect.ValueOf(HubLogIndexResp),
		"HubLogInfo":                        reflect.ValueOf(HubLogInfo),
		"IDPath":                            reflect.ValueOf(IDPath),
		"IDPathWithUpdate":                  reflect.ValueOf(IDPathWithUpdate),
		"OtaFirmwareDeviceCancelReq":        reflect.ValueOf(OtaFirmwareDeviceCancelReq),
		"OtaFirmwareDeviceConfirmReq":       reflect.ValueOf(OtaFirmwareDeviceConfirmReq),
		"OtaFirmwareDeviceIndexReq":         reflect.ValueOf(OtaFirmwareDeviceIndexReq),
		"OtaFirmwareDeviceIndexResp":        reflect.ValueOf(OtaFirmwareDeviceIndexResp),
		"OtaFirmwareDeviceInfo":             reflect.ValueOf(OtaFirmwareDeviceInfo),
		"OtaFirmwareDeviceRetryReq":         reflect.ValueOf(OtaFirmwareDeviceRetryReq),
		"OtaFirmwareFile":                   reflect.ValueOf(OtaFirmwareFile),
		"OtaFirmwareFileIndexReq":           reflect.ValueOf(OtaFirmwareFileIndexReq),
		"OtaFirmwareFileIndexResp":          reflect.ValueOf(OtaFirmwareFileIndexResp),
		"OtaFirmwareFileInfo":               reflect.ValueOf(OtaFirmwareFileInfo),
		"OtaFirmwareFileReq":                reflect.ValueOf(OtaFirmwareFileReq),
		"OtaFirmwareFileResp":               reflect.ValueOf(OtaFirmwareFileResp),
		"OtaFirmwareInfo":                   reflect.ValueOf(OtaFirmwareInfo),
		"OtaFirmwareInfoCreateReq":          reflect.ValueOf(OtaFirmwareInfoCreateReq),
		"OtaFirmwareInfoIndexReq":           reflect.ValueOf(OtaFirmwareInfoIndexReq),
		"OtaFirmwareInfoIndexResp":          reflect.ValueOf(OtaFirmwareInfoIndexResp),
		"OtaFirmwareInfoUpdateReq":          reflect.ValueOf(OtaFirmwareInfoUpdateReq),
		"OtaFirmwareJobIndexReq":            reflect.ValueOf(OtaFirmwareJobIndexReq),
		"OtaFirmwareJobIndexResp":           reflect.ValueOf(OtaFirmwareJobIndexResp),
		"OtaFirmwareJobInfo":                reflect.ValueOf(OtaFirmwareJobInfo),
		"OtaJobByDeviceIndexReq":            reflect.ValueOf(OtaJobByDeviceIndexReq),
		"OtaJobDynamicInfo":                 reflect.ValueOf(OtaJobDynamicInfo),
		"OtaJobStaticInfo":                  reflect.ValueOf(OtaJobStaticInfo),
		"OtaModuleInfo":                     reflect.ValueOf(OtaModuleInfo),
		"OtaModuleInfoIndexReq":             reflect.ValueOf(OtaModuleInfoIndexReq),
		"OtaModuleInfoIndexResp":            reflect.ValueOf(OtaModuleInfoIndexResp),
		"PageInfo":                          reflect.ValueOf(PageInfo),
		"PageInfo_OrderBy":                  reflect.ValueOf(PageInfo_OrderBy),
		"Point":                             reflect.ValueOf(Point),
		"ProductCategory":                   reflect.ValueOf(ProductCategory),
		"ProductCategoryIndexReq":           reflect.ValueOf(ProductCategoryIndexReq),
		"ProductCategoryIndexResp":          reflect.ValueOf(ProductCategoryIndexResp),
		"ProductCategorySchemaIndexReq":     reflect.ValueOf(ProductCategorySchemaIndexReq),
		"ProductCategorySchemaIndexResp":    reflect.ValueOf(ProductCategorySchemaIndexResp),
		"ProductCategorySchemaMultiSaveReq": reflect.ValueOf(ProductCategorySchemaMultiSaveReq),
		"ProductCustom":                     reflect.ValueOf(ProductCustom),
		"ProductCustomReadReq":              reflect.ValueOf(ProductCustomReadReq),
		"ProductCustomUi":                   reflect.ValueOf(ProductCustomUi),
		"ProductInfo":                       reflect.ValueOf(ProductInfo),
		"ProductInfoDeleteReq":              reflect.ValueOf(ProductInfoDeleteReq),
		"ProductInfoIndexReq":               reflect.ValueOf(ProductInfoIndexReq),
		"ProductInfoIndexResp":              reflect.ValueOf(ProductInfoIndexResp),
		"ProductInfoReadReq":                reflect.ValueOf(ProductInfoReadReq),
		"ProductInitReq":                    reflect.ValueOf(ProductInitReq),
		"ProductRemoteConfig":               reflect.ValueOf(ProductRemoteConfig),
		"ProductSchemaCreateReq":            reflect.ValueOf(ProductSchemaCreateReq),
		"ProductSchemaDeleteReq":            reflect.ValueOf(ProductSchemaDeleteReq),
		"ProductSchemaIndexReq":             reflect.ValueOf(ProductSchemaIndexReq),
		"ProductSchemaIndexResp":            reflect.ValueOf(ProductSchemaIndexResp),
		"ProductSchemaInfo":                 reflect.ValueOf(ProductSchemaInfo),
		"ProductSchemaMultiCreateReq":       reflect.ValueOf(ProductSchemaMultiCreateReq),
		"ProductSchemaTslImportReq":         reflect.ValueOf(ProductSchemaTslImportReq),
		"ProductSchemaTslReadReq":           reflect.ValueOf(ProductSchemaTslReadReq),
		"ProductSchemaTslReadResp":          reflect.ValueOf(ProductSchemaTslReadResp),
		"ProductSchemaUpdateReq":            reflect.ValueOf(ProductSchemaUpdateReq),
		"PropertyControlMultiSendReq":       reflect.ValueOf(PropertyControlMultiSendReq),
		"PropertyControlMultiSendResp":      reflect.ValueOf(PropertyControlMultiSendResp),
		"PropertyControlSendMsg":            reflect.ValueOf(PropertyControlSendMsg),
		"PropertyControlSendReq":            reflect.ValueOf(PropertyControlSendReq),
		"PropertyControlSendResp":           reflect.ValueOf(PropertyControlSendResp),
		"PropertyGetReportMultiSendReq":     reflect.ValueOf(PropertyGetReportMultiSendReq),
		"PropertyGetReportMultiSendResp":    reflect.ValueOf(PropertyGetReportMultiSendResp),
		"PropertyGetReportSendMsg":          reflect.ValueOf(PropertyGetReportSendMsg),
		"PropertyGetReportSendReq":          reflect.ValueOf(PropertyGetReportSendReq),
		"PropertyGetReportSendResp":         reflect.ValueOf(PropertyGetReportSendResp),
		"PropertyLogIndexReq":               reflect.ValueOf(PropertyLogIndexReq),
		"PropertyLogIndexResp":              reflect.ValueOf(PropertyLogIndexResp),
		"PropertyLogInfo":                   reflect.ValueOf(PropertyLogInfo),
		"PropertyLogLatestIndexReq":         reflect.ValueOf(PropertyLogLatestIndexReq),
		"ProtocolConfigField":               reflect.ValueOf(ProtocolConfigField),
		"ProtocolConfigInfo":                reflect.ValueOf(ProtocolConfigInfo),
		"ProtocolInfo":                      reflect.ValueOf(ProtocolInfo),
		"ProtocolInfoIndexReq":              reflect.ValueOf(ProtocolInfoIndexReq),
		"ProtocolInfoIndexResp":             reflect.ValueOf(ProtocolInfoIndexResp),
		"ProtocolService":                   reflect.ValueOf(ProtocolService),
		"ProtocolServiceIndexReq":           reflect.ValueOf(ProtocolServiceIndexReq),
		"ProtocolServiceIndexResp":          reflect.ValueOf(ProtocolServiceIndexResp),
		"RemoteConfigCreateReq":             reflect.ValueOf(RemoteConfigCreateReq),
		"RemoteConfigIndexReq":              reflect.ValueOf(RemoteConfigIndexReq),
		"RemoteConfigIndexResp":             reflect.ValueOf(RemoteConfigIndexResp),
		"RemoteConfigLastReadReq":           reflect.ValueOf(RemoteConfigLastReadReq),
		"RemoteConfigLastReadResp":          reflect.ValueOf(RemoteConfigLastReadResp),
		"RemoteConfigPushAllReq":            reflect.ValueOf(RemoteConfigPushAllReq),
		"RespReadReq":                       reflect.ValueOf(RespReadReq),
		"RootCheckReq":                      reflect.ValueOf(RootCheckReq),
		"SdkLogIndexReq":                    reflect.ValueOf(SdkLogIndexReq),
		"SdkLogIndexResp":                   reflect.ValueOf(SdkLogIndexResp),
		"SdkLogInfo":                        reflect.ValueOf(SdkLogInfo),
		"SendLogIndexReq":                   reflect.ValueOf(SendLogIndexReq),
		"SendLogIndexResp":                  reflect.ValueOf(SendLogIndexResp),
		"SendLogInfo":                       reflect.ValueOf(SendLogInfo),
		"SendMsgReq":                        reflect.ValueOf(SendMsgReq),
		"SendMsgResp":                       reflect.ValueOf(SendMsgResp),
		"SendOption":                        reflect.ValueOf(SendOption),
		"ShadowIndex":                       reflect.ValueOf(ShadowIndex),
		"ShadowIndexResp":                   reflect.ValueOf(ShadowIndexResp),
		"SharePerm":                         reflect.ValueOf(SharePerm),
		"StatusLogIndexReq":                 reflect.ValueOf(StatusLogIndexReq),
		"StatusLogIndexResp":                reflect.ValueOf(StatusLogIndexResp),
		"StatusLogInfo":                     reflect.ValueOf(StatusLogInfo),
		"TimeRange":                         reflect.ValueOf(TimeRange),
		"UserDeviceCollectSave":             reflect.ValueOf(UserDeviceCollectSave),
		"UserDeviceShareIndexReq":           reflect.ValueOf(UserDeviceShareIndexReq),
		"UserDeviceShareIndexResp":          reflect.ValueOf(UserDeviceShareIndexResp),
		"UserDeviceShareInfo":               reflect.ValueOf(UserDeviceShareInfo),
		"UserDeviceShareMultiAcceptReq":     reflect.ValueOf(UserDeviceShareMultiAcceptReq),
		"UserDeviceShareMultiDeleteReq":     reflect.ValueOf(UserDeviceShareMultiDeleteReq),
		"UserDeviceShareMultiInfo":          reflect.ValueOf(UserDeviceShareMultiInfo),
		"UserDeviceShareMultiToken":         reflect.ValueOf(UserDeviceShareMultiToken),
		"UserDeviceShareReadReq":            reflect.ValueOf(UserDeviceShareReadReq),
		"WithID":                            reflect.ValueOf(WithID),
		"WithIDChildren":                    reflect.ValueOf(WithIDChildren),
		"WithIDCode":                        reflect.ValueOf(WithIDCode),
		"WithProfile":                       reflect.ValueOf(WithProfile),
	}
}
