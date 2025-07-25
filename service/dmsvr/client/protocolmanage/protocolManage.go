// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: dm.proto

package protocolmanage

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AbnormalLogIndexReq               = dm.AbnormalLogIndexReq
	AbnormalLogIndexResp              = dm.AbnormalLogIndexResp
	AbnormalLogInfo                   = dm.AbnormalLogInfo
	ActionRespReq                     = dm.ActionRespReq
	ActionSendReq                     = dm.ActionSendReq
	ActionSendResp                    = dm.ActionSendResp
	CommonSchemaCreateReq             = dm.CommonSchemaCreateReq
	CommonSchemaExportReq             = dm.CommonSchemaExportReq
	CommonSchemaExportResp            = dm.CommonSchemaExportResp
	CommonSchemaImportReq             = dm.CommonSchemaImportReq
	CommonSchemaIndexReq              = dm.CommonSchemaIndexReq
	CommonSchemaIndexResp             = dm.CommonSchemaIndexResp
	CommonSchemaInfo                  = dm.CommonSchemaInfo
	CommonSchemaUpdateReq             = dm.CommonSchemaUpdateReq
	CompareInt64                      = dm.CompareInt64
	CompareString                     = dm.CompareString
	CustomTopic                       = dm.CustomTopic
	DevInit                           = dm.DevInit
	DeviceBindTokenInfo               = dm.DeviceBindTokenInfo
	DeviceBindTokenReadReq            = dm.DeviceBindTokenReadReq
	DeviceCore                        = dm.DeviceCore
	DeviceCountInfo                   = dm.DeviceCountInfo
	DeviceCountReq                    = dm.DeviceCountReq
	DeviceCountResp                   = dm.DeviceCountResp
	DeviceData                        = dm.DeviceData
	DeviceError                       = dm.DeviceError
	DeviceGatewayBindDevice           = dm.DeviceGatewayBindDevice
	DeviceGatewayIndexReq             = dm.DeviceGatewayIndexReq
	DeviceGatewayIndexResp            = dm.DeviceGatewayIndexResp
	DeviceGatewayMultiCreateReq       = dm.DeviceGatewayMultiCreateReq
	DeviceGatewayMultiSaveReq         = dm.DeviceGatewayMultiSaveReq
	DeviceGatewaySign                 = dm.DeviceGatewaySign
	DeviceGroupMultiSaveReq           = dm.DeviceGroupMultiSaveReq
	DeviceInfo                        = dm.DeviceInfo
	DeviceInfoBindReq                 = dm.DeviceInfoBindReq
	DeviceInfoCanBindReq              = dm.DeviceInfoCanBindReq
	DeviceInfoCount                   = dm.DeviceInfoCount
	DeviceInfoCountReq                = dm.DeviceInfoCountReq
	DeviceInfoDeleteReq               = dm.DeviceInfoDeleteReq
	DeviceInfoIndexReq                = dm.DeviceInfoIndexReq
	DeviceInfoIndexResp               = dm.DeviceInfoIndexResp
	DeviceInfoMultiBindReq            = dm.DeviceInfoMultiBindReq
	DeviceInfoMultiBindResp           = dm.DeviceInfoMultiBindResp
	DeviceInfoMultiUpdateReq          = dm.DeviceInfoMultiUpdateReq
	DeviceInfoReadReq                 = dm.DeviceInfoReadReq
	DeviceInfoUnbindReq               = dm.DeviceInfoUnbindReq
	DeviceModuleVersion               = dm.DeviceModuleVersion
	DeviceModuleVersionIndexReq       = dm.DeviceModuleVersionIndexReq
	DeviceModuleVersionIndexResp      = dm.DeviceModuleVersionIndexResp
	DeviceModuleVersionReadReq        = dm.DeviceModuleVersionReadReq
	DeviceMoveReq                     = dm.DeviceMoveReq
	DeviceOnlineMultiFix              = dm.DeviceOnlineMultiFix
	DeviceOnlineMultiFixReq           = dm.DeviceOnlineMultiFixReq
	DeviceProfile                     = dm.DeviceProfile
	DeviceProfileIndexReq             = dm.DeviceProfileIndexReq
	DeviceProfileIndexResp            = dm.DeviceProfileIndexResp
	DeviceProfileReadReq              = dm.DeviceProfileReadReq
	DeviceResetReq                    = dm.DeviceResetReq
	DeviceSchema                      = dm.DeviceSchema
	DeviceSchemaIndexReq              = dm.DeviceSchemaIndexReq
	DeviceSchemaIndexResp             = dm.DeviceSchemaIndexResp
	DeviceSchemaMultiCreateReq        = dm.DeviceSchemaMultiCreateReq
	DeviceSchemaMultiDeleteReq        = dm.DeviceSchemaMultiDeleteReq
	DeviceSchemaTslReadReq            = dm.DeviceSchemaTslReadReq
	DeviceSchemaTslReadResp           = dm.DeviceSchemaTslReadResp
	DeviceShareInfo                   = dm.DeviceShareInfo
	DeviceTransferReq                 = dm.DeviceTransferReq
	DeviceTypeCountReq                = dm.DeviceTypeCountReq
	DeviceTypeCountResp               = dm.DeviceTypeCountResp
	EdgeSendReq                       = dm.EdgeSendReq
	EdgeSendResp                      = dm.EdgeSendResp
	Empty                             = dm.Empty
	EventLogIndexReq                  = dm.EventLogIndexReq
	EventLogIndexResp                 = dm.EventLogIndexResp
	EventLogInfo                      = dm.EventLogInfo
	FileCore                          = dm.FileCore
	Firmware                          = dm.Firmware
	FirmwareFile                      = dm.FirmwareFile
	FirmwareInfo                      = dm.FirmwareInfo
	FirmwareInfoDeleteReq             = dm.FirmwareInfoDeleteReq
	FirmwareInfoDeleteResp            = dm.FirmwareInfoDeleteResp
	FirmwareInfoIndexReq              = dm.FirmwareInfoIndexReq
	FirmwareInfoIndexResp             = dm.FirmwareInfoIndexResp
	FirmwareInfoReadReq               = dm.FirmwareInfoReadReq
	FirmwareInfoReadResp              = dm.FirmwareInfoReadResp
	FirmwareResp                      = dm.FirmwareResp
	GatewayCanBindIndexReq            = dm.GatewayCanBindIndexReq
	GatewayCanBindIndexResp           = dm.GatewayCanBindIndexResp
	GatewayGetFoundReq                = dm.GatewayGetFoundReq
	GatewayNotifyBindSendReq          = dm.GatewayNotifyBindSendReq
	GroupCore                         = dm.GroupCore
	GroupDeviceMultiDeleteReq         = dm.GroupDeviceMultiDeleteReq
	GroupDeviceMultiSaveReq           = dm.GroupDeviceMultiSaveReq
	GroupInfo                         = dm.GroupInfo
	GroupInfoCreateReq                = dm.GroupInfoCreateReq
	GroupInfoIndexReq                 = dm.GroupInfoIndexReq
	GroupInfoIndexResp                = dm.GroupInfoIndexResp
	GroupInfoMultiCreateReq           = dm.GroupInfoMultiCreateReq
	GroupInfoReadReq                  = dm.GroupInfoReadReq
	GroupInfoUpdateReq                = dm.GroupInfoUpdateReq
	HubLogIndexReq                    = dm.HubLogIndexReq
	HubLogIndexResp                   = dm.HubLogIndexResp
	HubLogInfo                        = dm.HubLogInfo
	IDPath                            = dm.IDPath
	IDPathWithUpdate                  = dm.IDPathWithUpdate
	IDsInfo                           = dm.IDsInfo
	ImportResp                        = dm.ImportResp
	OtaFirmwareDeviceCancelReq        = dm.OtaFirmwareDeviceCancelReq
	OtaFirmwareDeviceConfirmReq       = dm.OtaFirmwareDeviceConfirmReq
	OtaFirmwareDeviceIndexReq         = dm.OtaFirmwareDeviceIndexReq
	OtaFirmwareDeviceIndexResp        = dm.OtaFirmwareDeviceIndexResp
	OtaFirmwareDeviceInfo             = dm.OtaFirmwareDeviceInfo
	OtaFirmwareDeviceRetryReq         = dm.OtaFirmwareDeviceRetryReq
	OtaFirmwareFile                   = dm.OtaFirmwareFile
	OtaFirmwareFileIndexReq           = dm.OtaFirmwareFileIndexReq
	OtaFirmwareFileIndexResp          = dm.OtaFirmwareFileIndexResp
	OtaFirmwareFileInfo               = dm.OtaFirmwareFileInfo
	OtaFirmwareFileReq                = dm.OtaFirmwareFileReq
	OtaFirmwareFileResp               = dm.OtaFirmwareFileResp
	OtaFirmwareInfo                   = dm.OtaFirmwareInfo
	OtaFirmwareInfoCreateReq          = dm.OtaFirmwareInfoCreateReq
	OtaFirmwareInfoIndexReq           = dm.OtaFirmwareInfoIndexReq
	OtaFirmwareInfoIndexResp          = dm.OtaFirmwareInfoIndexResp
	OtaFirmwareInfoUpdateReq          = dm.OtaFirmwareInfoUpdateReq
	OtaFirmwareJobIndexReq            = dm.OtaFirmwareJobIndexReq
	OtaFirmwareJobIndexResp           = dm.OtaFirmwareJobIndexResp
	OtaFirmwareJobInfo                = dm.OtaFirmwareJobInfo
	OtaJobByDeviceIndexReq            = dm.OtaJobByDeviceIndexReq
	OtaJobDynamicInfo                 = dm.OtaJobDynamicInfo
	OtaJobStaticInfo                  = dm.OtaJobStaticInfo
	OtaModuleInfo                     = dm.OtaModuleInfo
	OtaModuleInfoIndexReq             = dm.OtaModuleInfoIndexReq
	OtaModuleInfoIndexResp            = dm.OtaModuleInfoIndexResp
	PageInfo                          = dm.PageInfo
	PageInfo_OrderBy                  = dm.PageInfo_OrderBy
	Point                             = dm.Point
	ProductCategory                   = dm.ProductCategory
	ProductCategoryExportReq          = dm.ProductCategoryExportReq
	ProductCategoryExportResp         = dm.ProductCategoryExportResp
	ProductCategoryImportReq          = dm.ProductCategoryImportReq
	ProductCategoryIndexReq           = dm.ProductCategoryIndexReq
	ProductCategoryIndexResp          = dm.ProductCategoryIndexResp
	ProductCategorySchemaIndexReq     = dm.ProductCategorySchemaIndexReq
	ProductCategorySchemaIndexResp    = dm.ProductCategorySchemaIndexResp
	ProductCategorySchemaMultiSaveReq = dm.ProductCategorySchemaMultiSaveReq
	ProductConfig                     = dm.ProductConfig
	ProductCustom                     = dm.ProductCustom
	ProductCustomReadReq              = dm.ProductCustomReadReq
	ProductCustomUi                   = dm.ProductCustomUi
	ProductInfo                       = dm.ProductInfo
	ProductInfoDeleteReq              = dm.ProductInfoDeleteReq
	ProductInfoExportReq              = dm.ProductInfoExportReq
	ProductInfoExportResp             = dm.ProductInfoExportResp
	ProductInfoImportReq              = dm.ProductInfoImportReq
	ProductInfoIndexReq               = dm.ProductInfoIndexReq
	ProductInfoIndexResp              = dm.ProductInfoIndexResp
	ProductInfoReadReq                = dm.ProductInfoReadReq
	ProductInitReq                    = dm.ProductInitReq
	ProductRemoteConfig               = dm.ProductRemoteConfig
	ProductSchemaCreateReq            = dm.ProductSchemaCreateReq
	ProductSchemaDeleteReq            = dm.ProductSchemaDeleteReq
	ProductSchemaIndexReq             = dm.ProductSchemaIndexReq
	ProductSchemaIndexResp            = dm.ProductSchemaIndexResp
	ProductSchemaInfo                 = dm.ProductSchemaInfo
	ProductSchemaMultiCreateReq       = dm.ProductSchemaMultiCreateReq
	ProductSchemaTslImportReq         = dm.ProductSchemaTslImportReq
	ProductSchemaTslReadReq           = dm.ProductSchemaTslReadReq
	ProductSchemaTslReadResp          = dm.ProductSchemaTslReadResp
	ProductSchemaUpdateReq            = dm.ProductSchemaUpdateReq
	PropertyAgg                       = dm.PropertyAgg
	PropertyAggIndexReq               = dm.PropertyAggIndexReq
	PropertyAggIndexResp              = dm.PropertyAggIndexResp
	PropertyAggResp                   = dm.PropertyAggResp
	PropertyAggRespDataDetail         = dm.PropertyAggRespDataDetail
	PropertyAggRespDetail             = dm.PropertyAggRespDetail
	PropertyControlMultiSendReq       = dm.PropertyControlMultiSendReq
	PropertyControlMultiSendResp      = dm.PropertyControlMultiSendResp
	PropertyControlSendMsg            = dm.PropertyControlSendMsg
	PropertyControlSendReq            = dm.PropertyControlSendReq
	PropertyControlSendResp           = dm.PropertyControlSendResp
	PropertyGetReportMultiSendReq     = dm.PropertyGetReportMultiSendReq
	PropertyGetReportMultiSendResp    = dm.PropertyGetReportMultiSendResp
	PropertyGetReportSendMsg          = dm.PropertyGetReportSendMsg
	PropertyGetReportSendReq          = dm.PropertyGetReportSendReq
	PropertyGetReportSendResp         = dm.PropertyGetReportSendResp
	PropertyLogIndexReq               = dm.PropertyLogIndexReq
	PropertyLogIndexResp              = dm.PropertyLogIndexResp
	PropertyLogInfo                   = dm.PropertyLogInfo
	PropertyLogLatestIndex2Req        = dm.PropertyLogLatestIndex2Req
	PropertyLogLatestIndexReq         = dm.PropertyLogLatestIndexReq
	ProtocolConfigField               = dm.ProtocolConfigField
	ProtocolConfigInfo                = dm.ProtocolConfigInfo
	ProtocolInfo                      = dm.ProtocolInfo
	ProtocolInfoIndexReq              = dm.ProtocolInfoIndexReq
	ProtocolInfoIndexResp             = dm.ProtocolInfoIndexResp
	ProtocolScript                    = dm.ProtocolScript
	ProtocolScriptDebugReq            = dm.ProtocolScriptDebugReq
	ProtocolScriptDebugResp           = dm.ProtocolScriptDebugResp
	ProtocolScriptDevice              = dm.ProtocolScriptDevice
	ProtocolScriptDeviceIndexReq      = dm.ProtocolScriptDeviceIndexReq
	ProtocolScriptDeviceIndexResp     = dm.ProtocolScriptDeviceIndexResp
	ProtocolScriptExportReq           = dm.ProtocolScriptExportReq
	ProtocolScriptExportResp          = dm.ProtocolScriptExportResp
	ProtocolScriptImportReq           = dm.ProtocolScriptImportReq
	ProtocolScriptIndexReq            = dm.ProtocolScriptIndexReq
	ProtocolScriptIndexResp           = dm.ProtocolScriptIndexResp
	ProtocolService                   = dm.ProtocolService
	ProtocolServiceIndexReq           = dm.ProtocolServiceIndexReq
	ProtocolServiceIndexResp          = dm.ProtocolServiceIndexResp
	PublishMsg                        = dm.PublishMsg
	RemoteConfigCreateReq             = dm.RemoteConfigCreateReq
	RemoteConfigIndexReq              = dm.RemoteConfigIndexReq
	RemoteConfigIndexResp             = dm.RemoteConfigIndexResp
	RemoteConfigLastReadReq           = dm.RemoteConfigLastReadReq
	RemoteConfigLastReadResp          = dm.RemoteConfigLastReadResp
	RemoteConfigPushAllReq            = dm.RemoteConfigPushAllReq
	RespReadReq                       = dm.RespReadReq
	RootCheckReq                      = dm.RootCheckReq
	SdkLogIndexReq                    = dm.SdkLogIndexReq
	SdkLogIndexResp                   = dm.SdkLogIndexResp
	SdkLogInfo                        = dm.SdkLogInfo
	SendLogIndexReq                   = dm.SendLogIndexReq
	SendLogIndexResp                  = dm.SendLogIndexResp
	SendLogInfo                       = dm.SendLogInfo
	SendMsgReq                        = dm.SendMsgReq
	SendMsgResp                       = dm.SendMsgResp
	SendOption                        = dm.SendOption
	ShadowIndex                       = dm.ShadowIndex
	ShadowIndexResp                   = dm.ShadowIndexResp
	SharePerm                         = dm.SharePerm
	StatusLogIndexReq                 = dm.StatusLogIndexReq
	StatusLogIndexResp                = dm.StatusLogIndexResp
	StatusLogInfo                     = dm.StatusLogInfo
	TimeRange                         = dm.TimeRange
	UserDeviceCollectSave             = dm.UserDeviceCollectSave
	UserDeviceShareIndexReq           = dm.UserDeviceShareIndexReq
	UserDeviceShareIndexResp          = dm.UserDeviceShareIndexResp
	UserDeviceShareInfo               = dm.UserDeviceShareInfo
	UserDeviceShareMultiAcceptReq     = dm.UserDeviceShareMultiAcceptReq
	UserDeviceShareMultiDeleteReq     = dm.UserDeviceShareMultiDeleteReq
	UserDeviceShareMultiInfo          = dm.UserDeviceShareMultiInfo
	UserDeviceShareMultiToken         = dm.UserDeviceShareMultiToken
	UserDeviceShareReadReq            = dm.UserDeviceShareReadReq
	WithID                            = dm.WithID
	WithIDChildren                    = dm.WithIDChildren
	WithIDCode                        = dm.WithIDCode
	WithProfile                       = dm.WithProfile

	ProtocolManage interface {
		// 协议列表
		ProtocolInfoIndex(ctx context.Context, in *ProtocolInfoIndexReq, opts ...grpc.CallOption) (*ProtocolInfoIndexResp, error)
		// 协议详情
		ProtocolInfoRead(ctx context.Context, in *WithIDCode, opts ...grpc.CallOption) (*ProtocolInfo, error)
		// 协议创建
		ProtocolInfoCreate(ctx context.Context, in *ProtocolInfo, opts ...grpc.CallOption) (*WithID, error)
		// 协议更新
		ProtocolInfoUpdate(ctx context.Context, in *ProtocolInfo, opts ...grpc.CallOption) (*Empty, error)
		// 协议删除
		ProtocolInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
		// 更新服务状态,只给服务调用
		ProtocolServiceUpdate(ctx context.Context, in *ProtocolService, opts ...grpc.CallOption) (*Empty, error)
		ProtocolServiceDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
		ProtocolServiceIndex(ctx context.Context, in *ProtocolServiceIndexReq, opts ...grpc.CallOption) (*ProtocolServiceIndexResp, error)
		// 协议列表
		ProtocolScriptIndex(ctx context.Context, in *ProtocolScriptIndexReq, opts ...grpc.CallOption) (*ProtocolScriptIndexResp, error)
		// 协议详情
		ProtocolScriptRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*ProtocolScript, error)
		// 协议创建
		ProtocolScriptCreate(ctx context.Context, in *ProtocolScript, opts ...grpc.CallOption) (*WithID, error)
		// 协议更新
		ProtocolScriptUpdate(ctx context.Context, in *ProtocolScript, opts ...grpc.CallOption) (*Empty, error)
		// 协议删除
		ProtocolScriptDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
		ProtocolScriptDebug(ctx context.Context, in *ProtocolScriptDebugReq, opts ...grpc.CallOption) (*ProtocolScriptDebugResp, error)
		// 协议列表
		ProtocolScriptDeviceIndex(ctx context.Context, in *ProtocolScriptDeviceIndexReq, opts ...grpc.CallOption) (*ProtocolScriptDeviceIndexResp, error)
		// 协议详情
		ProtocolScriptDeviceRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*ProtocolScriptDevice, error)
		// 协议创建
		ProtocolScriptDeviceCreate(ctx context.Context, in *ProtocolScriptDevice, opts ...grpc.CallOption) (*WithID, error)
		// 协议更新
		ProtocolScriptDeviceUpdate(ctx context.Context, in *ProtocolScriptDevice, opts ...grpc.CallOption) (*Empty, error)
		// 协议删除
		ProtocolScriptDeviceDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error)
		ProtocolScriptMultiImport(ctx context.Context, in *ProtocolScriptImportReq, opts ...grpc.CallOption) (*ImportResp, error)
		ProtocolScriptMultiExport(ctx context.Context, in *ProtocolScriptExportReq, opts ...grpc.CallOption) (*ProtocolScriptExportResp, error)
	}

	defaultProtocolManage struct {
		cli zrpc.Client
	}

	directProtocolManage struct {
		svcCtx *svc.ServiceContext
		svr    dm.ProtocolManageServer
	}
)

func NewProtocolManage(cli zrpc.Client) ProtocolManage {
	return &defaultProtocolManage{
		cli: cli,
	}
}

func NewDirectProtocolManage(svcCtx *svc.ServiceContext, svr dm.ProtocolManageServer) ProtocolManage {
	return &directProtocolManage{
		svr:    svr,
		svcCtx: svcCtx,
	}
}

// 协议列表
func (m *defaultProtocolManage) ProtocolInfoIndex(ctx context.Context, in *ProtocolInfoIndexReq, opts ...grpc.CallOption) (*ProtocolInfoIndexResp, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolInfoIndex(ctx, in, opts...)
}

// 协议列表
func (d *directProtocolManage) ProtocolInfoIndex(ctx context.Context, in *ProtocolInfoIndexReq, opts ...grpc.CallOption) (*ProtocolInfoIndexResp, error) {
	return d.svr.ProtocolInfoIndex(ctx, in)
}

// 协议详情
func (m *defaultProtocolManage) ProtocolInfoRead(ctx context.Context, in *WithIDCode, opts ...grpc.CallOption) (*ProtocolInfo, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolInfoRead(ctx, in, opts...)
}

// 协议详情
func (d *directProtocolManage) ProtocolInfoRead(ctx context.Context, in *WithIDCode, opts ...grpc.CallOption) (*ProtocolInfo, error) {
	return d.svr.ProtocolInfoRead(ctx, in)
}

// 协议创建
func (m *defaultProtocolManage) ProtocolInfoCreate(ctx context.Context, in *ProtocolInfo, opts ...grpc.CallOption) (*WithID, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolInfoCreate(ctx, in, opts...)
}

// 协议创建
func (d *directProtocolManage) ProtocolInfoCreate(ctx context.Context, in *ProtocolInfo, opts ...grpc.CallOption) (*WithID, error) {
	return d.svr.ProtocolInfoCreate(ctx, in)
}

// 协议更新
func (m *defaultProtocolManage) ProtocolInfoUpdate(ctx context.Context, in *ProtocolInfo, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolInfoUpdate(ctx, in, opts...)
}

// 协议更新
func (d *directProtocolManage) ProtocolInfoUpdate(ctx context.Context, in *ProtocolInfo, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolInfoUpdate(ctx, in)
}

// 协议删除
func (m *defaultProtocolManage) ProtocolInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolInfoDelete(ctx, in, opts...)
}

// 协议删除
func (d *directProtocolManage) ProtocolInfoDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolInfoDelete(ctx, in)
}

// 更新服务状态,只给服务调用
func (m *defaultProtocolManage) ProtocolServiceUpdate(ctx context.Context, in *ProtocolService, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolServiceUpdate(ctx, in, opts...)
}

// 更新服务状态,只给服务调用
func (d *directProtocolManage) ProtocolServiceUpdate(ctx context.Context, in *ProtocolService, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolServiceUpdate(ctx, in)
}

func (m *defaultProtocolManage) ProtocolServiceDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolServiceDelete(ctx, in, opts...)
}

func (d *directProtocolManage) ProtocolServiceDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolServiceDelete(ctx, in)
}

func (m *defaultProtocolManage) ProtocolServiceIndex(ctx context.Context, in *ProtocolServiceIndexReq, opts ...grpc.CallOption) (*ProtocolServiceIndexResp, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolServiceIndex(ctx, in, opts...)
}

func (d *directProtocolManage) ProtocolServiceIndex(ctx context.Context, in *ProtocolServiceIndexReq, opts ...grpc.CallOption) (*ProtocolServiceIndexResp, error) {
	return d.svr.ProtocolServiceIndex(ctx, in)
}

// 协议列表
func (m *defaultProtocolManage) ProtocolScriptIndex(ctx context.Context, in *ProtocolScriptIndexReq, opts ...grpc.CallOption) (*ProtocolScriptIndexResp, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptIndex(ctx, in, opts...)
}

// 协议列表
func (d *directProtocolManage) ProtocolScriptIndex(ctx context.Context, in *ProtocolScriptIndexReq, opts ...grpc.CallOption) (*ProtocolScriptIndexResp, error) {
	return d.svr.ProtocolScriptIndex(ctx, in)
}

// 协议详情
func (m *defaultProtocolManage) ProtocolScriptRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*ProtocolScript, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptRead(ctx, in, opts...)
}

// 协议详情
func (d *directProtocolManage) ProtocolScriptRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*ProtocolScript, error) {
	return d.svr.ProtocolScriptRead(ctx, in)
}

// 协议创建
func (m *defaultProtocolManage) ProtocolScriptCreate(ctx context.Context, in *ProtocolScript, opts ...grpc.CallOption) (*WithID, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptCreate(ctx, in, opts...)
}

// 协议创建
func (d *directProtocolManage) ProtocolScriptCreate(ctx context.Context, in *ProtocolScript, opts ...grpc.CallOption) (*WithID, error) {
	return d.svr.ProtocolScriptCreate(ctx, in)
}

// 协议更新
func (m *defaultProtocolManage) ProtocolScriptUpdate(ctx context.Context, in *ProtocolScript, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptUpdate(ctx, in, opts...)
}

// 协议更新
func (d *directProtocolManage) ProtocolScriptUpdate(ctx context.Context, in *ProtocolScript, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolScriptUpdate(ctx, in)
}

// 协议删除
func (m *defaultProtocolManage) ProtocolScriptDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptDelete(ctx, in, opts...)
}

// 协议删除
func (d *directProtocolManage) ProtocolScriptDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolScriptDelete(ctx, in)
}

func (m *defaultProtocolManage) ProtocolScriptDebug(ctx context.Context, in *ProtocolScriptDebugReq, opts ...grpc.CallOption) (*ProtocolScriptDebugResp, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptDebug(ctx, in, opts...)
}

func (d *directProtocolManage) ProtocolScriptDebug(ctx context.Context, in *ProtocolScriptDebugReq, opts ...grpc.CallOption) (*ProtocolScriptDebugResp, error) {
	return d.svr.ProtocolScriptDebug(ctx, in)
}

// 协议列表
func (m *defaultProtocolManage) ProtocolScriptDeviceIndex(ctx context.Context, in *ProtocolScriptDeviceIndexReq, opts ...grpc.CallOption) (*ProtocolScriptDeviceIndexResp, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptDeviceIndex(ctx, in, opts...)
}

// 协议列表
func (d *directProtocolManage) ProtocolScriptDeviceIndex(ctx context.Context, in *ProtocolScriptDeviceIndexReq, opts ...grpc.CallOption) (*ProtocolScriptDeviceIndexResp, error) {
	return d.svr.ProtocolScriptDeviceIndex(ctx, in)
}

// 协议详情
func (m *defaultProtocolManage) ProtocolScriptDeviceRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*ProtocolScriptDevice, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptDeviceRead(ctx, in, opts...)
}

// 协议详情
func (d *directProtocolManage) ProtocolScriptDeviceRead(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*ProtocolScriptDevice, error) {
	return d.svr.ProtocolScriptDeviceRead(ctx, in)
}

// 协议创建
func (m *defaultProtocolManage) ProtocolScriptDeviceCreate(ctx context.Context, in *ProtocolScriptDevice, opts ...grpc.CallOption) (*WithID, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptDeviceCreate(ctx, in, opts...)
}

// 协议创建
func (d *directProtocolManage) ProtocolScriptDeviceCreate(ctx context.Context, in *ProtocolScriptDevice, opts ...grpc.CallOption) (*WithID, error) {
	return d.svr.ProtocolScriptDeviceCreate(ctx, in)
}

// 协议更新
func (m *defaultProtocolManage) ProtocolScriptDeviceUpdate(ctx context.Context, in *ProtocolScriptDevice, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptDeviceUpdate(ctx, in, opts...)
}

// 协议更新
func (d *directProtocolManage) ProtocolScriptDeviceUpdate(ctx context.Context, in *ProtocolScriptDevice, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolScriptDeviceUpdate(ctx, in)
}

// 协议删除
func (m *defaultProtocolManage) ProtocolScriptDeviceDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptDeviceDelete(ctx, in, opts...)
}

// 协议删除
func (d *directProtocolManage) ProtocolScriptDeviceDelete(ctx context.Context, in *WithID, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.ProtocolScriptDeviceDelete(ctx, in)
}

func (m *defaultProtocolManage) ProtocolScriptMultiImport(ctx context.Context, in *ProtocolScriptImportReq, opts ...grpc.CallOption) (*ImportResp, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptMultiImport(ctx, in, opts...)
}

func (d *directProtocolManage) ProtocolScriptMultiImport(ctx context.Context, in *ProtocolScriptImportReq, opts ...grpc.CallOption) (*ImportResp, error) {
	return d.svr.ProtocolScriptMultiImport(ctx, in)
}

func (m *defaultProtocolManage) ProtocolScriptMultiExport(ctx context.Context, in *ProtocolScriptExportReq, opts ...grpc.CallOption) (*ProtocolScriptExportResp, error) {
	client := dm.NewProtocolManageClient(m.cli.Conn())
	return client.ProtocolScriptMultiExport(ctx, in, opts...)
}

func (d *directProtocolManage) ProtocolScriptMultiExport(ctx context.Context, in *ProtocolScriptExportReq, opts ...grpc.CallOption) (*ProtocolScriptExportResp, error) {
	return d.svr.ProtocolScriptMultiExport(ctx, in)
}
