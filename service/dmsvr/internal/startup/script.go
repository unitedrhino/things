package startup

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/things/service/dmsvr/client/deviceinteract"
	"gitee.com/unitedrhino/things/service/dmsvr/client/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/otamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/client/productmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	deviceinteractServer "gitee.com/unitedrhino/things/service/dmsvr/internal/server/deviceinteract"
	devicemanageServer "gitee.com/unitedrhino/things/service/dmsvr/internal/server/devicemanage"
	otamanageServer "gitee.com/unitedrhino/things/service/dmsvr/internal/server/otamanage"
	productmanageServer "gitee.com/unitedrhino/things/service/dmsvr/internal/server/productmanage"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgExt"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgGateway"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgSdkLog"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
)

func ScriptInit(svcCtx *svc.ServiceContext) {
	ScriptLoad(svcCtx)
	svcCtx.ScriptTrans.AddSymbol("dm/dm", dmSymbolInit(svcCtx))
	svcCtx.ScriptTrans.AddSymbol("schema/schema", schemaSymbolInit(svcCtx))
	svcCtx.ScriptTrans.AddSymbol("deviceMsg/deviceMsg", deviceMsgSymbolInit(svcCtx))

	return
}

func ScriptLoad(svcCtx *svc.ServiceContext) {
	svcCtx.ScriptTrans.SetLoad(func(ctx context.Context, trans *protocol.ScriptTrans) error {
		sds, err := relationDB.NewProtocolScriptDeviceRepo(ctx).FindByFilter(ctx, relationDB.ProtocolScriptDeviceFilter{
			WithScript: true, Status: def.True}, nil)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		var (
			ProductCache = make(map[protocol.TriggerDir]map[protocol.TriggerTimer]map[string]map[devices.MsgHandle]map[string]protocol.ScriptInfos)
			DeviceCache  = make(map[protocol.TriggerDir]map[protocol.TriggerTimer]map[devices.Core]map[devices.MsgHandle]map[string]protocol.ScriptInfos)
		)
		for _, sd := range sds {
			if sd.Script == nil || sd.Script.Status == def.False {
				continue
			}
			si := protocol.ScriptInfo{
				Name:       sd.Script.Name,
				Priority:   sd.Priority,
				ScriptLang: sd.Script.ScriptLang,
				Script:     sd.Script.Script,
			}
			switch sd.TriggerSrc {
			case protocol.TriggerSrcDevice:
				dev := devices.Core{ProductID: sd.ProductID, DeviceName: sd.DeviceName}
				_, ok := DeviceCache[sd.Script.TriggerDir]
				if !ok {
					DeviceCache[sd.Script.TriggerDir] = map[protocol.TriggerTimer]map[devices.Core]map[devices.MsgHandle]map[string]protocol.ScriptInfos{
						sd.Script.TriggerTimer: {dev: {sd.Script.TriggerHandle: {sd.Script.TriggerType: protocol.ScriptInfos{si}}}},
					}
					continue
				}
				_, ok = DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer]
				if !ok {
					DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer] = map[devices.Core]map[devices.MsgHandle]map[string]protocol.ScriptInfos{
						dev: {sd.Script.TriggerHandle: {sd.Script.TriggerType: protocol.ScriptInfos{si}}},
					}
					continue
				}
				_, ok = DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][dev]
				if !ok {
					DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][dev] = map[devices.MsgHandle]map[string]protocol.ScriptInfos{
						sd.Script.TriggerHandle: {sd.Script.TriggerType: protocol.ScriptInfos{si}},
					}
					continue
				}
				_, ok = DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][dev][sd.Script.TriggerHandle]
				if !ok {
					DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][dev][sd.Script.TriggerHandle] = map[string]protocol.ScriptInfos{
						sd.Script.TriggerType: protocol.ScriptInfos{si},
					}
					continue
				}
				DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][dev][sd.Script.TriggerHandle][sd.Script.TriggerType] =
					append(DeviceCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][dev][sd.Script.TriggerHandle][sd.Script.TriggerType], si)
			case protocol.TriggerSrcProduct:
				_, ok := ProductCache[sd.Script.TriggerDir]
				if !ok {
					ProductCache[sd.Script.TriggerDir] = map[protocol.TriggerTimer]map[string]map[devices.MsgHandle]map[string]protocol.ScriptInfos{
						sd.Script.TriggerTimer: {sd.ProductID: {sd.Script.TriggerHandle: {sd.Script.TriggerType: protocol.ScriptInfos{si}}}},
					}
					continue
				}
				_, ok = ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer]
				if !ok {
					ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer] = map[string]map[devices.MsgHandle]map[string]protocol.ScriptInfos{
						sd.ProductID: {sd.Script.TriggerHandle: {sd.Script.TriggerType: protocol.ScriptInfos{si}}},
					}
					continue
				}
				_, ok = ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][sd.ProductID]
				if !ok {
					ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][sd.ProductID] = map[devices.MsgHandle]map[string]protocol.ScriptInfos{
						sd.Script.TriggerHandle: {sd.Script.TriggerType: protocol.ScriptInfos{si}},
					}
					continue
				}
				_, ok = ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][sd.ProductID][sd.Script.TriggerHandle]
				if !ok {
					ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][sd.ProductID][sd.Script.TriggerHandle] = map[string]protocol.ScriptInfos{
						sd.Script.TriggerType: protocol.ScriptInfos{si},
					}
					continue
				}
				ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][sd.ProductID][sd.Script.TriggerHandle][sd.Script.TriggerType] =
					append(ProductCache[sd.Script.TriggerDir][sd.Script.TriggerTimer][sd.ProductID][sd.Script.TriggerHandle][sd.Script.TriggerType], si)
			}
		}
		if len(DeviceCache) > 0 {
			up := DeviceCache[protocol.TriggerDirUp]
			down := DeviceCache[protocol.TriggerDirDown]
			if len(up) > 0 {
				before := up[protocol.TriggerTimerBefore]
				if len(before) > 0 {
					func() {
						trans.DeviceUpBeforeMutex.Lock()
						defer trans.DeviceUpBeforeMutex.Unlock()
						trans.DeviceUpBeforeCache = before
					}()
				}
			}
			if len(down) > 0 {
				before := down[protocol.TriggerTimerBefore]
				if len(before) > 0 {
					func() {
						trans.DeviceDownBeforeMutex.Lock()
						defer trans.DeviceDownBeforeMutex.Unlock()
						trans.DeviceDownBeforeCache = before
					}()
				}
			}
		}
		if len(ProductCache) > 0 {
			up := ProductCache[protocol.TriggerDirUp]
			down := ProductCache[protocol.TriggerDirDown]
			if len(up) > 0 {
				before := up[protocol.TriggerTimerBefore]
				if len(before) > 0 {
					func() {
						trans.ProductUpBeforeMutex.Lock()
						defer trans.ProductUpBeforeMutex.Unlock()
						trans.ProductUpBeforeCache = before
					}()
				}
			}
			if len(down) > 0 {
				before := down[protocol.TriggerTimerBefore]
				if len(before) > 0 {
					func() {
						trans.ProductDownBeforeMutex.Lock()
						defer trans.ProductDownBeforeMutex.Unlock()
						trans.ProductDownBeforeCache = before
					}()
				}
			}
		}
		return nil
	})
}

func schemaSymbolInit(svcCtx *svc.ServiceContext) map[string]reflect.Value {
	return map[string]reflect.Value{
		"ModelSimple": reflect.ValueOf((*schema.ModelSimple)(nil)),
		"Model":       reflect.ValueOf((*schema.Model)(nil)),
	}
}

func deviceMsgSymbolInit(svcCtx *svc.ServiceContext) map[string]reflect.Value {
	return map[string]reflect.Value{
		"PublishMsg": reflect.ValueOf((*deviceMsg.PublishMsg)(nil)),
		"CommonMsg":  reflect.ValueOf((*deviceMsg.CommonMsg)(nil)),
		"TimeParams": reflect.ValueOf((*deviceMsg.TimeParams)(nil)),

		"SysConfig":      reflect.ValueOf((*deviceMsg.SysConfig)(nil)),
		"ThingReq":       reflect.ValueOf((*msgThing.Req)(nil)),
		"ThingSubDevice": reflect.ValueOf((*msgThing.SubDevice)(nil)),
		"ThingResp":      reflect.ValueOf((*msgThing.Resp)(nil)),

		"SdkLogReq":    reflect.ValueOf((*msgSdkLog.Req)(nil)),
		"SdkLogSdkLog": reflect.ValueOf((*msgSdkLog.SdkLog)(nil)),

		"OtaReq":           reflect.ValueOf((*msgOta.Req)(nil)),
		"OtaProcess":       reflect.ValueOf((*msgOta.Process)(nil)),
		"OtaParams":        reflect.ValueOf((*msgOta.Params)(nil)),
		"OtaProcessParams": reflect.ValueOf((*msgOta.ProcessParams)(nil)),
		"OtaUpgrade":       reflect.ValueOf((*msgOta.Upgrade)(nil)),

		"GatewayMsg":      reflect.ValueOf((*msgGateway.Msg)(nil)),
		"GatewayPayload":  reflect.ValueOf((*msgGateway.GatewayPayload)(nil)),
		"GatewayRegister": reflect.ValueOf((*msgGateway.Register)(nil)),
		"GatewayDevice":   reflect.ValueOf((*msgGateway.Device)(nil)),

		"ExtReq":          reflect.ValueOf((*msgExt.Req)(nil)),
		"ExtRegisterReq":  reflect.ValueOf((*msgExt.RegisterReq)(nil)),
		"ExtResp":         reflect.ValueOf((*msgExt.Resp)(nil)),
		"ExtRespRegister": reflect.ValueOf((*msgExt.RespRegister)(nil)),
	}
}

func dmSymbolInit(svcCtx *svc.ServiceContext) map[string]reflect.Value {
	return map[string]reflect.Value{
		"ProductGet": reflect.ValueOf(func(ctx context.Context, productID string) (*dm.ProductInfo, error) {
			return svcCtx.ProductCache.GetData(ctx, productID)
		}),
		"DeviceGet": reflect.ValueOf(func(ctx context.Context, productID string, deviceName string) (*dm.DeviceInfo, error) {
			return svcCtx.DeviceCache.GetData(ctx, devices.Core{ProductID: productID, DeviceName: deviceName})
		}),
		"SchemaGet": reflect.ValueOf(func(ctx context.Context, productID string, deviceName string) (*schema.Model, error) {
			return svcCtx.DeviceSchemaRepo.GetData(ctx, devices.Core{ProductID: productID, DeviceName: deviceName})
		}),
		"DeviceInteract":                    reflect.ValueOf(deviceinteract.NewDirectDeviceInteract(svcCtx, deviceinteractServer.NewDeviceInteractServer(svcCtx))),
		"DeviceManage":                      reflect.ValueOf(devicemanage.NewDirectDeviceManage(svcCtx, devicemanageServer.NewDeviceManageServer(svcCtx))),
		"ProductManage":                     reflect.ValueOf(productmanage.NewDirectProductManage(svcCtx, productmanageServer.NewProductManageServer(svcCtx))),
		"OtaManage":                         reflect.ValueOf(otamanage.NewDirectOtaManage(svcCtx, otamanageServer.NewOtaManageServer(svcCtx))),
		"AbnormalLogIndexReq":               reflect.ValueOf((*dm.AbnormalLogIndexReq)(nil)),
		"AbnormalLogIndexResp":              reflect.ValueOf((*dm.AbnormalLogIndexResp)(nil)),
		"AbnormalLogInfo":                   reflect.ValueOf((*dm.AbnormalLogInfo)(nil)),
		"ActionRespReq":                     reflect.ValueOf((*dm.ActionRespReq)(nil)),
		"ActionSendReq":                     reflect.ValueOf((*dm.ActionSendReq)(nil)),
		"ActionSendResp":                    reflect.ValueOf((*dm.ActionSendResp)(nil)),
		"CommonSchemaCreateReq":             reflect.ValueOf((*dm.CommonSchemaCreateReq)(nil)),
		"CommonSchemaIndexReq":              reflect.ValueOf((*dm.CommonSchemaIndexReq)(nil)),
		"CommonSchemaIndexResp":             reflect.ValueOf((*dm.CommonSchemaIndexResp)(nil)),
		"CommonSchemaInfo":                  reflect.ValueOf((*dm.CommonSchemaInfo)(nil)),
		"CommonSchemaUpdateReq":             reflect.ValueOf((*dm.CommonSchemaUpdateReq)(nil)),
		"CompareInt64":                      reflect.ValueOf((*dm.CompareInt64)(nil)),
		"CompareString":                     reflect.ValueOf((*dm.CompareString)(nil)),
		"CustomTopic":                       reflect.ValueOf((*dm.CustomTopic)(nil)),
		"DeviceBindTokenInfo":               reflect.ValueOf((*dm.DeviceBindTokenInfo)(nil)),
		"DeviceBindTokenReadReq":            reflect.ValueOf((*dm.DeviceBindTokenReadReq)(nil)),
		"DeviceCore":                        reflect.ValueOf((*dm.DeviceCore)(nil)),
		"DeviceCountInfo":                   reflect.ValueOf((*dm.DeviceCountInfo)(nil)),
		"DeviceCountReq":                    reflect.ValueOf((*dm.DeviceCountReq)(nil)),
		"DeviceCountResp":                   reflect.ValueOf((*dm.DeviceCountResp)(nil)),
		"DeviceError":                       reflect.ValueOf((*dm.DeviceError)(nil)),
		"DeviceGatewayBindDevice":           reflect.ValueOf((*dm.DeviceGatewayBindDevice)(nil)),
		"DeviceGatewayIndexReq":             reflect.ValueOf((*dm.DeviceGatewayIndexReq)(nil)),
		"DeviceGatewayIndexResp":            reflect.ValueOf((*dm.DeviceGatewayIndexResp)(nil)),
		"DeviceGatewayMultiCreateReq":       reflect.ValueOf((*dm.DeviceGatewayMultiCreateReq)(nil)),
		"DeviceGatewayMultiSaveReq":         reflect.ValueOf((*dm.DeviceGatewayMultiSaveReq)(nil)),
		"DeviceGatewaySign":                 reflect.ValueOf((*dm.DeviceGatewaySign)(nil)),
		"DeviceGroupMultiSaveReq":           reflect.ValueOf((*dm.DeviceGroupMultiSaveReq)(nil)),
		"DeviceInfo":                        reflect.ValueOf((*dm.DeviceInfo)(nil)),
		"DeviceInfoBindReq":                 reflect.ValueOf((*dm.DeviceInfoBindReq)(nil)),
		"DeviceInfoCanBindReq":              reflect.ValueOf((*dm.DeviceInfoCanBindReq)(nil)),
		"DeviceInfoCount":                   reflect.ValueOf((*dm.DeviceInfoCount)(nil)),
		"DeviceInfoCountReq":                reflect.ValueOf((*dm.DeviceInfoCountReq)(nil)),
		"DeviceInfoDeleteReq":               reflect.ValueOf((*dm.DeviceInfoDeleteReq)(nil)),
		"DeviceInfoIndexReq":                reflect.ValueOf((*dm.DeviceInfoIndexReq)(nil)),
		"DeviceInfoIndexResp":               reflect.ValueOf((*dm.DeviceInfoIndexResp)(nil)),
		"DeviceInfoMultiBindReq":            reflect.ValueOf((*dm.DeviceInfoMultiBindReq)(nil)),
		"DeviceInfoMultiBindResp":           reflect.ValueOf((*dm.DeviceInfoMultiBindResp)(nil)),
		"DeviceInfoMultiUpdateReq":          reflect.ValueOf((*dm.DeviceInfoMultiUpdateReq)(nil)),
		"DeviceInfoReadReq":                 reflect.ValueOf((*dm.DeviceInfoReadReq)(nil)),
		"DeviceInfoUnbindReq":               reflect.ValueOf((*dm.DeviceInfoUnbindReq)(nil)),
		"DeviceModuleVersion":               reflect.ValueOf((*dm.DeviceModuleVersion)(nil)),
		"DeviceModuleVersionIndexReq":       reflect.ValueOf((*dm.DeviceModuleVersionIndexReq)(nil)),
		"DeviceModuleVersionIndexResp":      reflect.ValueOf((*dm.DeviceModuleVersionIndexResp)(nil)),
		"DeviceModuleVersionReadReq":        reflect.ValueOf((*dm.DeviceModuleVersionReadReq)(nil)),
		"DeviceMoveReq":                     reflect.ValueOf((*dm.DeviceMoveReq)(nil)),
		"DeviceOnlineMultiFix":              reflect.ValueOf((*dm.DeviceOnlineMultiFix)(nil)),
		"DeviceOnlineMultiFixReq":           reflect.ValueOf((*dm.DeviceOnlineMultiFixReq)(nil)),
		"DeviceProfile":                     reflect.ValueOf((*dm.DeviceProfile)(nil)),
		"DeviceProfileIndexReq":             reflect.ValueOf((*dm.DeviceProfileIndexReq)(nil)),
		"DeviceProfileIndexResp":            reflect.ValueOf((*dm.DeviceProfileIndexResp)(nil)),
		"DeviceProfileReadReq":              reflect.ValueOf((*dm.DeviceProfileReadReq)(nil)),
		"DeviceResetReq":                    reflect.ValueOf((*dm.DeviceResetReq)(nil)),
		"DeviceSchema":                      reflect.ValueOf((*dm.DeviceSchema)(nil)),
		"DeviceSchemaIndexReq":              reflect.ValueOf((*dm.DeviceSchemaIndexReq)(nil)),
		"DeviceSchemaIndexResp":             reflect.ValueOf((*dm.DeviceSchemaIndexResp)(nil)),
		"DeviceSchemaMultiCreateReq":        reflect.ValueOf((*dm.DeviceSchemaMultiCreateReq)(nil)),
		"DeviceSchemaMultiDeleteReq":        reflect.ValueOf((*dm.DeviceSchemaMultiDeleteReq)(nil)),
		"DeviceSchemaTslReadReq":            reflect.ValueOf((*dm.DeviceSchemaTslReadReq)(nil)),
		"DeviceSchemaTslReadResp":           reflect.ValueOf((*dm.DeviceSchemaTslReadResp)(nil)),
		"DeviceShareInfo":                   reflect.ValueOf((*dm.DeviceShareInfo)(nil)),
		"DeviceTransferReq":                 reflect.ValueOf((*dm.DeviceTransferReq)(nil)),
		"DeviceTypeCountReq":                reflect.ValueOf((*dm.DeviceTypeCountReq)(nil)),
		"DeviceTypeCountResp":               reflect.ValueOf((*dm.DeviceTypeCountResp)(nil)),
		"EdgeSendReq":                       reflect.ValueOf((*dm.EdgeSendReq)(nil)),
		"EdgeSendResp":                      reflect.ValueOf((*dm.EdgeSendResp)(nil)),
		"Empty":                             reflect.ValueOf((*dm.Empty)(nil)),
		"EventLogIndexReq":                  reflect.ValueOf((*dm.EventLogIndexReq)(nil)),
		"EventLogIndexResp":                 reflect.ValueOf((*dm.EventLogIndexResp)(nil)),
		"EventLogInfo":                      reflect.ValueOf((*dm.EventLogInfo)(nil)),
		"Firmware":                          reflect.ValueOf((*dm.Firmware)(nil)),
		"FirmwareFile":                      reflect.ValueOf((*dm.FirmwareFile)(nil)),
		"FirmwareInfo":                      reflect.ValueOf((*dm.FirmwareInfo)(nil)),
		"FirmwareInfoDeleteReq":             reflect.ValueOf((*dm.FirmwareInfoDeleteReq)(nil)),
		"FirmwareInfoDeleteResp":            reflect.ValueOf((*dm.FirmwareInfoDeleteResp)(nil)),
		"FirmwareInfoIndexReq":              reflect.ValueOf((*dm.FirmwareInfoIndexReq)(nil)),
		"FirmwareInfoIndexResp":             reflect.ValueOf((*dm.FirmwareInfoIndexResp)(nil)),
		"FirmwareInfoReadReq":               reflect.ValueOf((*dm.FirmwareInfoReadReq)(nil)),
		"FirmwareInfoReadResp":              reflect.ValueOf((*dm.FirmwareInfoReadResp)(nil)),
		"FirmwareResp":                      reflect.ValueOf((*dm.FirmwareResp)(nil)),
		"GatewayCanBindIndexReq":            reflect.ValueOf((*dm.GatewayCanBindIndexReq)(nil)),
		"GatewayCanBindIndexResp":           reflect.ValueOf((*dm.GatewayCanBindIndexResp)(nil)),
		"GatewayGetFoundReq":                reflect.ValueOf((*dm.GatewayGetFoundReq)(nil)),
		"GatewayNotifyBindSendReq":          reflect.ValueOf((*dm.GatewayNotifyBindSendReq)(nil)),
		"GroupCore":                         reflect.ValueOf((*dm.GroupCore)(nil)),
		"GroupDeviceMultiDeleteReq":         reflect.ValueOf((*dm.GroupDeviceMultiDeleteReq)(nil)),
		"GroupDeviceMultiSaveReq":           reflect.ValueOf((*dm.GroupDeviceMultiSaveReq)(nil)),
		"GroupInfo":                         reflect.ValueOf((*dm.GroupInfo)(nil)),
		"GroupInfoCreateReq":                reflect.ValueOf((*dm.GroupInfoCreateReq)(nil)),
		"GroupInfoIndexReq":                 reflect.ValueOf((*dm.GroupInfoIndexReq)(nil)),
		"GroupInfoIndexResp":                reflect.ValueOf((*dm.GroupInfoIndexResp)(nil)),
		"GroupInfoMultiCreateReq":           reflect.ValueOf((*dm.GroupInfoMultiCreateReq)(nil)),
		"GroupInfoReadReq":                  reflect.ValueOf((*dm.GroupInfoReadReq)(nil)),
		"GroupInfoUpdateReq":                reflect.ValueOf((*dm.GroupInfoUpdateReq)(nil)),
		"HubLogIndexReq":                    reflect.ValueOf((*dm.HubLogIndexReq)(nil)),
		"HubLogIndexResp":                   reflect.ValueOf((*dm.HubLogIndexResp)(nil)),
		"HubLogInfo":                        reflect.ValueOf((*dm.HubLogInfo)(nil)),
		"IDPath":                            reflect.ValueOf((*dm.IDPath)(nil)),
		"IDPathWithUpdate":                  reflect.ValueOf((*dm.IDPathWithUpdate)(nil)),
		"OtaFirmwareDeviceCancelReq":        reflect.ValueOf((*dm.OtaFirmwareDeviceCancelReq)(nil)),
		"OtaFirmwareDeviceConfirmReq":       reflect.ValueOf((*dm.OtaFirmwareDeviceConfirmReq)(nil)),
		"OtaFirmwareDeviceIndexReq":         reflect.ValueOf((*dm.OtaFirmwareDeviceIndexReq)(nil)),
		"OtaFirmwareDeviceIndexResp":        reflect.ValueOf((*dm.OtaFirmwareDeviceIndexResp)(nil)),
		"OtaFirmwareDeviceInfo":             reflect.ValueOf((*dm.OtaFirmwareDeviceInfo)(nil)),
		"OtaFirmwareDeviceRetryReq":         reflect.ValueOf((*dm.OtaFirmwareDeviceRetryReq)(nil)),
		"OtaFirmwareFile":                   reflect.ValueOf((*dm.OtaFirmwareFile)(nil)),
		"OtaFirmwareFileIndexReq":           reflect.ValueOf((*dm.OtaFirmwareFileIndexReq)(nil)),
		"OtaFirmwareFileIndexResp":          reflect.ValueOf((*dm.OtaFirmwareFileIndexResp)(nil)),
		"OtaFirmwareFileInfo":               reflect.ValueOf((*dm.OtaFirmwareFileInfo)(nil)),
		"OtaFirmwareFileReq":                reflect.ValueOf((*dm.OtaFirmwareFileReq)(nil)),
		"OtaFirmwareFileResp":               reflect.ValueOf((*dm.OtaFirmwareFileResp)(nil)),
		"OtaFirmwareInfo":                   reflect.ValueOf((*dm.OtaFirmwareInfo)(nil)),
		"OtaFirmwareInfoCreateReq":          reflect.ValueOf((*dm.OtaFirmwareInfoCreateReq)(nil)),
		"OtaFirmwareInfoIndexReq":           reflect.ValueOf((*dm.OtaFirmwareInfoIndexReq)(nil)),
		"OtaFirmwareInfoIndexResp":          reflect.ValueOf((*dm.OtaFirmwareInfoIndexResp)(nil)),
		"OtaFirmwareInfoUpdateReq":          reflect.ValueOf((*dm.OtaFirmwareInfoUpdateReq)(nil)),
		"OtaFirmwareJobIndexReq":            reflect.ValueOf((*dm.OtaFirmwareJobIndexReq)(nil)),
		"OtaFirmwareJobIndexResp":           reflect.ValueOf((*dm.OtaFirmwareJobIndexResp)(nil)),
		"OtaFirmwareJobInfo":                reflect.ValueOf((*dm.OtaFirmwareJobInfo)(nil)),
		"OtaJobByDeviceIndexReq":            reflect.ValueOf((*dm.OtaJobByDeviceIndexReq)(nil)),
		"OtaJobDynamicInfo":                 reflect.ValueOf((*dm.OtaJobDynamicInfo)(nil)),
		"OtaJobStaticInfo":                  reflect.ValueOf((*dm.OtaJobStaticInfo)(nil)),
		"OtaModuleInfo":                     reflect.ValueOf((*dm.OtaModuleInfo)(nil)),
		"OtaModuleInfoIndexReq":             reflect.ValueOf((*dm.OtaModuleInfoIndexReq)(nil)),
		"OtaModuleInfoIndexResp":            reflect.ValueOf((*dm.OtaModuleInfoIndexResp)(nil)),
		"PageInfo":                          reflect.ValueOf((*dm.PageInfo)(nil)),
		"PageInfo_OrderBy":                  reflect.ValueOf((*dm.PageInfo_OrderBy)(nil)),
		"Point":                             reflect.ValueOf((*dm.Point)(nil)),
		"ProductCategory":                   reflect.ValueOf((*dm.ProductCategory)(nil)),
		"ProductCategoryIndexReq":           reflect.ValueOf((*dm.ProductCategoryIndexReq)(nil)),
		"ProductCategoryIndexResp":          reflect.ValueOf((*dm.ProductCategoryIndexResp)(nil)),
		"ProductCategorySchemaIndexReq":     reflect.ValueOf((*dm.ProductCategorySchemaIndexReq)(nil)),
		"ProductCategorySchemaIndexResp":    reflect.ValueOf((*dm.ProductCategorySchemaIndexResp)(nil)),
		"ProductCategorySchemaMultiSaveReq": reflect.ValueOf((*dm.ProductCategorySchemaMultiSaveReq)(nil)),
		"ProductCustom":                     reflect.ValueOf((*dm.ProductCustom)(nil)),
		"ProductCustomReadReq":              reflect.ValueOf((*dm.ProductCustomReadReq)(nil)),
		"ProductCustomUi":                   reflect.ValueOf((*dm.ProductCustomUi)(nil)),
		"ProductInfo":                       reflect.ValueOf((*dm.ProductInfo)(nil)),
		"ProductInfoDeleteReq":              reflect.ValueOf((*dm.ProductInfoDeleteReq)(nil)),
		"ProductInfoIndexReq":               reflect.ValueOf((*dm.ProductInfoIndexReq)(nil)),
		"ProductInfoIndexResp":              reflect.ValueOf((*dm.ProductInfoIndexResp)(nil)),
		"ProductInfoReadReq":                reflect.ValueOf((*dm.ProductInfoReadReq)(nil)),
		"ProductInitReq":                    reflect.ValueOf((*dm.ProductInitReq)(nil)),
		"ProductRemoteConfig":               reflect.ValueOf((*dm.ProductRemoteConfig)(nil)),
		"ProductSchemaCreateReq":            reflect.ValueOf((*dm.ProductSchemaCreateReq)(nil)),
		"ProductSchemaDeleteReq":            reflect.ValueOf((*dm.ProductSchemaDeleteReq)(nil)),
		"ProductSchemaIndexReq":             reflect.ValueOf((*dm.ProductSchemaIndexReq)(nil)),
		"ProductSchemaIndexResp":            reflect.ValueOf((*dm.ProductSchemaIndexResp)(nil)),
		"ProductSchemaInfo":                 reflect.ValueOf((*dm.ProductSchemaInfo)(nil)),
		"ProductSchemaMultiCreateReq":       reflect.ValueOf((*dm.ProductSchemaMultiCreateReq)(nil)),
		"ProductSchemaTslImportReq":         reflect.ValueOf((*dm.ProductSchemaTslImportReq)(nil)),
		"ProductSchemaTslReadReq":           reflect.ValueOf((*dm.ProductSchemaTslReadReq)(nil)),
		"ProductSchemaTslReadResp":          reflect.ValueOf((*dm.ProductSchemaTslReadResp)(nil)),
		"ProductSchemaUpdateReq":            reflect.ValueOf((*dm.ProductSchemaUpdateReq)(nil)),
		"PropertyControlMultiSendReq":       reflect.ValueOf((*dm.PropertyControlMultiSendReq)(nil)),
		"PropertyControlMultiSendResp":      reflect.ValueOf((*dm.PropertyControlMultiSendResp)(nil)),
		"PropertyControlSendMsg":            reflect.ValueOf((*dm.PropertyControlSendMsg)(nil)),
		"PropertyControlSendReq":            reflect.ValueOf((*dm.PropertyControlSendReq)(nil)),
		"PropertyControlSendResp":           reflect.ValueOf((*dm.PropertyControlSendResp)(nil)),
		"PropertyGetReportMultiSendReq":     reflect.ValueOf((*dm.PropertyGetReportMultiSendReq)(nil)),
		"PropertyGetReportMultiSendResp":    reflect.ValueOf((*dm.PropertyGetReportMultiSendResp)(nil)),
		"PropertyGetReportSendMsg":          reflect.ValueOf((*dm.PropertyGetReportSendMsg)(nil)),
		"PropertyGetReportSendReq":          reflect.ValueOf((*dm.PropertyGetReportSendReq)(nil)),
		"PropertyGetReportSendResp":         reflect.ValueOf((*dm.PropertyGetReportSendResp)(nil)),
		"PropertyLogIndexReq":               reflect.ValueOf((*dm.PropertyLogIndexReq)(nil)),
		"PropertyLogIndexResp":              reflect.ValueOf((*dm.PropertyLogIndexResp)(nil)),
		"PropertyLogInfo":                   reflect.ValueOf((*dm.PropertyLogInfo)(nil)),
		"PropertyLogLatestIndexReq":         reflect.ValueOf((*dm.PropertyLogLatestIndexReq)(nil)),
		"ProtocolConfigField":               reflect.ValueOf((*dm.ProtocolConfigField)(nil)),
		"ProtocolConfigInfo":                reflect.ValueOf((*dm.ProtocolConfigInfo)(nil)),
		"ProtocolInfo":                      reflect.ValueOf((*dm.ProtocolInfo)(nil)),
		"ProtocolInfoIndexReq":              reflect.ValueOf((*dm.ProtocolInfoIndexReq)(nil)),
		"ProtocolInfoIndexResp":             reflect.ValueOf((*dm.ProtocolInfoIndexResp)(nil)),
		"ProtocolService":                   reflect.ValueOf((*dm.ProtocolService)(nil)),
		"ProtocolServiceIndexReq":           reflect.ValueOf((*dm.ProtocolServiceIndexReq)(nil)),
		"ProtocolServiceIndexResp":          reflect.ValueOf((*dm.ProtocolServiceIndexResp)(nil)),
		"RemoteConfigCreateReq":             reflect.ValueOf((*dm.RemoteConfigCreateReq)(nil)),
		"RemoteConfigIndexReq":              reflect.ValueOf((*dm.RemoteConfigIndexReq)(nil)),
		"RemoteConfigIndexResp":             reflect.ValueOf((*dm.RemoteConfigIndexResp)(nil)),
		"RemoteConfigLastReadReq":           reflect.ValueOf((*dm.RemoteConfigLastReadReq)(nil)),
		"RemoteConfigLastReadResp":          reflect.ValueOf((*dm.RemoteConfigLastReadResp)(nil)),
		"RemoteConfigPushAllReq":            reflect.ValueOf((*dm.RemoteConfigPushAllReq)(nil)),
		"RespReadReq":                       reflect.ValueOf((*dm.RespReadReq)(nil)),
		"RootCheckReq":                      reflect.ValueOf((*dm.RootCheckReq)(nil)),
		"SdkLogIndexReq":                    reflect.ValueOf((*dm.SdkLogIndexReq)(nil)),
		"SdkLogIndexResp":                   reflect.ValueOf((*dm.SdkLogIndexResp)(nil)),
		"SdkLogInfo":                        reflect.ValueOf((*dm.SdkLogInfo)(nil)),
		"SendLogIndexReq":                   reflect.ValueOf((*dm.SendLogIndexReq)(nil)),
		"SendLogIndexResp":                  reflect.ValueOf((*dm.SendLogIndexResp)(nil)),
		"SendLogInfo":                       reflect.ValueOf((*dm.SendLogInfo)(nil)),
		"SendMsgReq":                        reflect.ValueOf((*dm.SendMsgReq)(nil)),
		"SendMsgResp":                       reflect.ValueOf((*dm.SendMsgResp)(nil)),
		"SendOption":                        reflect.ValueOf((*dm.SendOption)(nil)),
		"ShadowIndex":                       reflect.ValueOf((*dm.ShadowIndex)(nil)),
		"ShadowIndexResp":                   reflect.ValueOf((*dm.ShadowIndexResp)(nil)),
		"SharePerm":                         reflect.ValueOf((*dm.SharePerm)(nil)),
		"StatusLogIndexReq":                 reflect.ValueOf((*dm.StatusLogIndexReq)(nil)),
		"StatusLogIndexResp":                reflect.ValueOf((*dm.StatusLogIndexResp)(nil)),
		"StatusLogInfo":                     reflect.ValueOf((*dm.StatusLogInfo)(nil)),
		"TimeRange":                         reflect.ValueOf((*dm.TimeRange)(nil)),
		"UserDeviceCollectSave":             reflect.ValueOf((*dm.UserDeviceCollectSave)(nil)),
		"UserDeviceShareIndexReq":           reflect.ValueOf((*dm.UserDeviceShareIndexReq)(nil)),
		"UserDeviceShareIndexResp":          reflect.ValueOf((*dm.UserDeviceShareIndexResp)(nil)),
		"UserDeviceShareInfo":               reflect.ValueOf((*dm.UserDeviceShareInfo)(nil)),
		"UserDeviceShareMultiAcceptReq":     reflect.ValueOf((*dm.UserDeviceShareMultiAcceptReq)(nil)),
		"UserDeviceShareMultiDeleteReq":     reflect.ValueOf((*dm.UserDeviceShareMultiDeleteReq)(nil)),
		"UserDeviceShareMultiInfo":          reflect.ValueOf((*dm.UserDeviceShareMultiInfo)(nil)),
		"UserDeviceShareMultiToken":         reflect.ValueOf((*dm.UserDeviceShareMultiToken)(nil)),
		"UserDeviceShareReadReq":            reflect.ValueOf((*dm.UserDeviceShareReadReq)(nil)),
		"WithID":                            reflect.ValueOf((*dm.WithID)(nil)),
		"WithIDChildren":                    reflect.ValueOf((*dm.WithIDChildren)(nil)),
		"WithIDCode":                        reflect.ValueOf((*dm.WithIDCode)(nil)),
		"WithProfile":                       reflect.ValueOf((*dm.WithProfile)(nil)),
	}
}
