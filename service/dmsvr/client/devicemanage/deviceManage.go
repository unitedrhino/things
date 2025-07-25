// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: dm.proto

package devicemanage

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

	DeviceManage interface {
		// 鉴定是否是root账号(提供给mqtt broker)
		RootCheck(ctx context.Context, in *RootCheckReq, opts ...grpc.CallOption) (*Empty, error)
		// 新增设备
		DeviceInfoCreate(ctx context.Context, in *DeviceInfo, opts ...grpc.CallOption) (*Empty, error)
		// 更新设备
		DeviceInfoUpdate(ctx context.Context, in *DeviceInfo, opts ...grpc.CallOption) (*Empty, error)
		DeviceOnlineMultiFix(ctx context.Context, in *DeviceOnlineMultiFixReq, opts ...grpc.CallOption) (*Empty, error)
		// 删除设备
		DeviceInfoDelete(ctx context.Context, in *DeviceInfoDeleteReq, opts ...grpc.CallOption) (*Empty, error)
		// 获取设备信息列表
		DeviceInfoIndex(ctx context.Context, in *DeviceInfoIndexReq, opts ...grpc.CallOption) (*DeviceInfoIndexResp, error)
		// 批量更新设备状态
		DeviceInfoMultiUpdate(ctx context.Context, in *DeviceInfoMultiUpdateReq, opts ...grpc.CallOption) (*Empty, error)
		// 获取设备信息详情
		DeviceInfoRead(ctx context.Context, in *DeviceInfoReadReq, opts ...grpc.CallOption) (*DeviceInfo, error)
		DeviceInfoBind(ctx context.Context, in *DeviceInfoBindReq, opts ...grpc.CallOption) (*Empty, error)
		DeviceBindTokenRead(ctx context.Context, in *DeviceBindTokenReadReq, opts ...grpc.CallOption) (*DeviceBindTokenInfo, error)
		DeviceBindTokenCreate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DeviceBindTokenInfo, error)
		DeviceInfoMultiBind(ctx context.Context, in *DeviceInfoMultiBindReq, opts ...grpc.CallOption) (*DeviceInfoMultiBindResp, error)
		DeviceInfoCanBind(ctx context.Context, in *DeviceInfoCanBindReq, opts ...grpc.CallOption) (*Empty, error)
		DeviceInfoUnbind(ctx context.Context, in *DeviceInfoUnbindReq, opts ...grpc.CallOption) (*Empty, error)
		DeviceTransfer(ctx context.Context, in *DeviceTransferReq, opts ...grpc.CallOption) (*Empty, error)
		DeviceReset(ctx context.Context, in *DeviceResetReq, opts ...grpc.CallOption) (*Empty, error)
		DeviceMove(ctx context.Context, in *DeviceMoveReq, opts ...grpc.CallOption) (*Empty, error)
		DeviceModuleVersionRead(ctx context.Context, in *DeviceModuleVersionReadReq, opts ...grpc.CallOption) (*DeviceModuleVersion, error)
		DeviceModuleVersionIndex(ctx context.Context, in *DeviceModuleVersionIndexReq, opts ...grpc.CallOption) (*DeviceModuleVersionIndexResp, error)
		// 绑定网关下子设备设备
		DeviceGatewayMultiCreate(ctx context.Context, in *DeviceGatewayMultiCreateReq, opts ...grpc.CallOption) (*Empty, error)
		// 绑定网关下子设备设备
		DeviceGatewayMultiUpdate(ctx context.Context, in *DeviceGatewayMultiSaveReq, opts ...grpc.CallOption) (*Empty, error)
		// 获取绑定信息的设备信息列表
		DeviceGatewayIndex(ctx context.Context, in *DeviceGatewayIndexReq, opts ...grpc.CallOption) (*DeviceGatewayIndexResp, error)
		// 删除网关下子设备
		DeviceGatewayMultiDelete(ctx context.Context, in *DeviceGatewayMultiSaveReq, opts ...grpc.CallOption) (*Empty, error)
		// 设备计数
		DeviceInfoCount(ctx context.Context, in *DeviceInfoCountReq, opts ...grpc.CallOption) (*DeviceInfoCount, error)
		// 设备类型
		DeviceTypeCount(ctx context.Context, in *DeviceTypeCountReq, opts ...grpc.CallOption) (*DeviceTypeCountResp, error)
		DeviceCount(ctx context.Context, in *DeviceCountReq, opts ...grpc.CallOption) (*DeviceCountResp, error)
		DeviceProfileRead(ctx context.Context, in *DeviceProfileReadReq, opts ...grpc.CallOption) (*DeviceProfile, error)
		DeviceProfileDelete(ctx context.Context, in *DeviceProfileReadReq, opts ...grpc.CallOption) (*Empty, error)
		DeviceProfileUpdate(ctx context.Context, in *DeviceProfile, opts ...grpc.CallOption) (*Empty, error)
		DeviceProfileIndex(ctx context.Context, in *DeviceProfileIndexReq, opts ...grpc.CallOption) (*DeviceProfileIndexResp, error)
		// 更新设备物模型
		DeviceSchemaUpdate(ctx context.Context, in *DeviceSchema, opts ...grpc.CallOption) (*Empty, error)
		// 新增设备
		DeviceSchemaCreate(ctx context.Context, in *DeviceSchema, opts ...grpc.CallOption) (*Empty, error)
		// 批量新增物模型,只新增没有的,已有的不处理
		DeviceSchemaMultiCreate(ctx context.Context, in *DeviceSchemaMultiCreateReq, opts ...grpc.CallOption) (*Empty, error)
		// 删除设备物模型
		DeviceSchemaMultiDelete(ctx context.Context, in *DeviceSchemaMultiDeleteReq, opts ...grpc.CallOption) (*Empty, error)
		// 获取设备物模型列表
		DeviceSchemaIndex(ctx context.Context, in *DeviceSchemaIndexReq, opts ...grpc.CallOption) (*DeviceSchemaIndexResp, error)
		DeviceSchemaTslRead(ctx context.Context, in *DeviceSchemaTslReadReq, opts ...grpc.CallOption) (*DeviceSchemaTslReadResp, error)
		// 将设备加到多个分组中
		DeviceGroupMultiCreate(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error)
		// 更新设备所在分组
		DeviceGroupMultiUpdate(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error)
		// 删除设备所在分组
		DeviceGroupMultiDelete(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error)
	}

	defaultDeviceManage struct {
		cli zrpc.Client
	}

	directDeviceManage struct {
		svcCtx *svc.ServiceContext
		svr    dm.DeviceManageServer
	}
)

func NewDeviceManage(cli zrpc.Client) DeviceManage {
	return &defaultDeviceManage{
		cli: cli,
	}
}

func NewDirectDeviceManage(svcCtx *svc.ServiceContext, svr dm.DeviceManageServer) DeviceManage {
	return &directDeviceManage{
		svr:    svr,
		svcCtx: svcCtx,
	}
}

// 鉴定是否是root账号(提供给mqtt broker)
func (m *defaultDeviceManage) RootCheck(ctx context.Context, in *RootCheckReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.RootCheck(ctx, in, opts...)
}

// 鉴定是否是root账号(提供给mqtt broker)
func (d *directDeviceManage) RootCheck(ctx context.Context, in *RootCheckReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.RootCheck(ctx, in)
}

// 新增设备
func (m *defaultDeviceManage) DeviceInfoCreate(ctx context.Context, in *DeviceInfo, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoCreate(ctx, in, opts...)
}

// 新增设备
func (d *directDeviceManage) DeviceInfoCreate(ctx context.Context, in *DeviceInfo, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceInfoCreate(ctx, in)
}

// 更新设备
func (m *defaultDeviceManage) DeviceInfoUpdate(ctx context.Context, in *DeviceInfo, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoUpdate(ctx, in, opts...)
}

// 更新设备
func (d *directDeviceManage) DeviceInfoUpdate(ctx context.Context, in *DeviceInfo, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceInfoUpdate(ctx, in)
}

func (m *defaultDeviceManage) DeviceOnlineMultiFix(ctx context.Context, in *DeviceOnlineMultiFixReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceOnlineMultiFix(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceOnlineMultiFix(ctx context.Context, in *DeviceOnlineMultiFixReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceOnlineMultiFix(ctx, in)
}

// 删除设备
func (m *defaultDeviceManage) DeviceInfoDelete(ctx context.Context, in *DeviceInfoDeleteReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoDelete(ctx, in, opts...)
}

// 删除设备
func (d *directDeviceManage) DeviceInfoDelete(ctx context.Context, in *DeviceInfoDeleteReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceInfoDelete(ctx, in)
}

// 获取设备信息列表
func (m *defaultDeviceManage) DeviceInfoIndex(ctx context.Context, in *DeviceInfoIndexReq, opts ...grpc.CallOption) (*DeviceInfoIndexResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoIndex(ctx, in, opts...)
}

// 获取设备信息列表
func (d *directDeviceManage) DeviceInfoIndex(ctx context.Context, in *DeviceInfoIndexReq, opts ...grpc.CallOption) (*DeviceInfoIndexResp, error) {
	return d.svr.DeviceInfoIndex(ctx, in)
}

// 批量更新设备状态
func (m *defaultDeviceManage) DeviceInfoMultiUpdate(ctx context.Context, in *DeviceInfoMultiUpdateReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoMultiUpdate(ctx, in, opts...)
}

// 批量更新设备状态
func (d *directDeviceManage) DeviceInfoMultiUpdate(ctx context.Context, in *DeviceInfoMultiUpdateReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceInfoMultiUpdate(ctx, in)
}

// 获取设备信息详情
func (m *defaultDeviceManage) DeviceInfoRead(ctx context.Context, in *DeviceInfoReadReq, opts ...grpc.CallOption) (*DeviceInfo, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoRead(ctx, in, opts...)
}

// 获取设备信息详情
func (d *directDeviceManage) DeviceInfoRead(ctx context.Context, in *DeviceInfoReadReq, opts ...grpc.CallOption) (*DeviceInfo, error) {
	return d.svr.DeviceInfoRead(ctx, in)
}

func (m *defaultDeviceManage) DeviceInfoBind(ctx context.Context, in *DeviceInfoBindReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoBind(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceInfoBind(ctx context.Context, in *DeviceInfoBindReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceInfoBind(ctx, in)
}

func (m *defaultDeviceManage) DeviceBindTokenRead(ctx context.Context, in *DeviceBindTokenReadReq, opts ...grpc.CallOption) (*DeviceBindTokenInfo, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceBindTokenRead(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceBindTokenRead(ctx context.Context, in *DeviceBindTokenReadReq, opts ...grpc.CallOption) (*DeviceBindTokenInfo, error) {
	return d.svr.DeviceBindTokenRead(ctx, in)
}

func (m *defaultDeviceManage) DeviceBindTokenCreate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DeviceBindTokenInfo, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceBindTokenCreate(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceBindTokenCreate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*DeviceBindTokenInfo, error) {
	return d.svr.DeviceBindTokenCreate(ctx, in)
}

func (m *defaultDeviceManage) DeviceInfoMultiBind(ctx context.Context, in *DeviceInfoMultiBindReq, opts ...grpc.CallOption) (*DeviceInfoMultiBindResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoMultiBind(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceInfoMultiBind(ctx context.Context, in *DeviceInfoMultiBindReq, opts ...grpc.CallOption) (*DeviceInfoMultiBindResp, error) {
	return d.svr.DeviceInfoMultiBind(ctx, in)
}

func (m *defaultDeviceManage) DeviceInfoCanBind(ctx context.Context, in *DeviceInfoCanBindReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoCanBind(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceInfoCanBind(ctx context.Context, in *DeviceInfoCanBindReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceInfoCanBind(ctx, in)
}

func (m *defaultDeviceManage) DeviceInfoUnbind(ctx context.Context, in *DeviceInfoUnbindReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoUnbind(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceInfoUnbind(ctx context.Context, in *DeviceInfoUnbindReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceInfoUnbind(ctx, in)
}

func (m *defaultDeviceManage) DeviceTransfer(ctx context.Context, in *DeviceTransferReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceTransfer(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceTransfer(ctx context.Context, in *DeviceTransferReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceTransfer(ctx, in)
}

func (m *defaultDeviceManage) DeviceReset(ctx context.Context, in *DeviceResetReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceReset(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceReset(ctx context.Context, in *DeviceResetReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceReset(ctx, in)
}

func (m *defaultDeviceManage) DeviceMove(ctx context.Context, in *DeviceMoveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceMove(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceMove(ctx context.Context, in *DeviceMoveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceMove(ctx, in)
}

func (m *defaultDeviceManage) DeviceModuleVersionRead(ctx context.Context, in *DeviceModuleVersionReadReq, opts ...grpc.CallOption) (*DeviceModuleVersion, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceModuleVersionRead(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceModuleVersionRead(ctx context.Context, in *DeviceModuleVersionReadReq, opts ...grpc.CallOption) (*DeviceModuleVersion, error) {
	return d.svr.DeviceModuleVersionRead(ctx, in)
}

func (m *defaultDeviceManage) DeviceModuleVersionIndex(ctx context.Context, in *DeviceModuleVersionIndexReq, opts ...grpc.CallOption) (*DeviceModuleVersionIndexResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceModuleVersionIndex(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceModuleVersionIndex(ctx context.Context, in *DeviceModuleVersionIndexReq, opts ...grpc.CallOption) (*DeviceModuleVersionIndexResp, error) {
	return d.svr.DeviceModuleVersionIndex(ctx, in)
}

// 绑定网关下子设备设备
func (m *defaultDeviceManage) DeviceGatewayMultiCreate(ctx context.Context, in *DeviceGatewayMultiCreateReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceGatewayMultiCreate(ctx, in, opts...)
}

// 绑定网关下子设备设备
func (d *directDeviceManage) DeviceGatewayMultiCreate(ctx context.Context, in *DeviceGatewayMultiCreateReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceGatewayMultiCreate(ctx, in)
}

// 绑定网关下子设备设备
func (m *defaultDeviceManage) DeviceGatewayMultiUpdate(ctx context.Context, in *DeviceGatewayMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceGatewayMultiUpdate(ctx, in, opts...)
}

// 绑定网关下子设备设备
func (d *directDeviceManage) DeviceGatewayMultiUpdate(ctx context.Context, in *DeviceGatewayMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceGatewayMultiUpdate(ctx, in)
}

// 获取绑定信息的设备信息列表
func (m *defaultDeviceManage) DeviceGatewayIndex(ctx context.Context, in *DeviceGatewayIndexReq, opts ...grpc.CallOption) (*DeviceGatewayIndexResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceGatewayIndex(ctx, in, opts...)
}

// 获取绑定信息的设备信息列表
func (d *directDeviceManage) DeviceGatewayIndex(ctx context.Context, in *DeviceGatewayIndexReq, opts ...grpc.CallOption) (*DeviceGatewayIndexResp, error) {
	return d.svr.DeviceGatewayIndex(ctx, in)
}

// 删除网关下子设备
func (m *defaultDeviceManage) DeviceGatewayMultiDelete(ctx context.Context, in *DeviceGatewayMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceGatewayMultiDelete(ctx, in, opts...)
}

// 删除网关下子设备
func (d *directDeviceManage) DeviceGatewayMultiDelete(ctx context.Context, in *DeviceGatewayMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceGatewayMultiDelete(ctx, in)
}

// 设备计数
func (m *defaultDeviceManage) DeviceInfoCount(ctx context.Context, in *DeviceInfoCountReq, opts ...grpc.CallOption) (*DeviceInfoCount, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceInfoCount(ctx, in, opts...)
}

// 设备计数
func (d *directDeviceManage) DeviceInfoCount(ctx context.Context, in *DeviceInfoCountReq, opts ...grpc.CallOption) (*DeviceInfoCount, error) {
	return d.svr.DeviceInfoCount(ctx, in)
}

// 设备类型
func (m *defaultDeviceManage) DeviceTypeCount(ctx context.Context, in *DeviceTypeCountReq, opts ...grpc.CallOption) (*DeviceTypeCountResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceTypeCount(ctx, in, opts...)
}

// 设备类型
func (d *directDeviceManage) DeviceTypeCount(ctx context.Context, in *DeviceTypeCountReq, opts ...grpc.CallOption) (*DeviceTypeCountResp, error) {
	return d.svr.DeviceTypeCount(ctx, in)
}

func (m *defaultDeviceManage) DeviceCount(ctx context.Context, in *DeviceCountReq, opts ...grpc.CallOption) (*DeviceCountResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceCount(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceCount(ctx context.Context, in *DeviceCountReq, opts ...grpc.CallOption) (*DeviceCountResp, error) {
	return d.svr.DeviceCount(ctx, in)
}

func (m *defaultDeviceManage) DeviceProfileRead(ctx context.Context, in *DeviceProfileReadReq, opts ...grpc.CallOption) (*DeviceProfile, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceProfileRead(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceProfileRead(ctx context.Context, in *DeviceProfileReadReq, opts ...grpc.CallOption) (*DeviceProfile, error) {
	return d.svr.DeviceProfileRead(ctx, in)
}

func (m *defaultDeviceManage) DeviceProfileDelete(ctx context.Context, in *DeviceProfileReadReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceProfileDelete(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceProfileDelete(ctx context.Context, in *DeviceProfileReadReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceProfileDelete(ctx, in)
}

func (m *defaultDeviceManage) DeviceProfileUpdate(ctx context.Context, in *DeviceProfile, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceProfileUpdate(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceProfileUpdate(ctx context.Context, in *DeviceProfile, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceProfileUpdate(ctx, in)
}

func (m *defaultDeviceManage) DeviceProfileIndex(ctx context.Context, in *DeviceProfileIndexReq, opts ...grpc.CallOption) (*DeviceProfileIndexResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceProfileIndex(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceProfileIndex(ctx context.Context, in *DeviceProfileIndexReq, opts ...grpc.CallOption) (*DeviceProfileIndexResp, error) {
	return d.svr.DeviceProfileIndex(ctx, in)
}

// 更新设备物模型
func (m *defaultDeviceManage) DeviceSchemaUpdate(ctx context.Context, in *DeviceSchema, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceSchemaUpdate(ctx, in, opts...)
}

// 更新设备物模型
func (d *directDeviceManage) DeviceSchemaUpdate(ctx context.Context, in *DeviceSchema, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceSchemaUpdate(ctx, in)
}

// 新增设备
func (m *defaultDeviceManage) DeviceSchemaCreate(ctx context.Context, in *DeviceSchema, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceSchemaCreate(ctx, in, opts...)
}

// 新增设备
func (d *directDeviceManage) DeviceSchemaCreate(ctx context.Context, in *DeviceSchema, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceSchemaCreate(ctx, in)
}

// 批量新增物模型,只新增没有的,已有的不处理
func (m *defaultDeviceManage) DeviceSchemaMultiCreate(ctx context.Context, in *DeviceSchemaMultiCreateReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceSchemaMultiCreate(ctx, in, opts...)
}

// 批量新增物模型,只新增没有的,已有的不处理
func (d *directDeviceManage) DeviceSchemaMultiCreate(ctx context.Context, in *DeviceSchemaMultiCreateReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceSchemaMultiCreate(ctx, in)
}

// 删除设备物模型
func (m *defaultDeviceManage) DeviceSchemaMultiDelete(ctx context.Context, in *DeviceSchemaMultiDeleteReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceSchemaMultiDelete(ctx, in, opts...)
}

// 删除设备物模型
func (d *directDeviceManage) DeviceSchemaMultiDelete(ctx context.Context, in *DeviceSchemaMultiDeleteReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceSchemaMultiDelete(ctx, in)
}

// 获取设备物模型列表
func (m *defaultDeviceManage) DeviceSchemaIndex(ctx context.Context, in *DeviceSchemaIndexReq, opts ...grpc.CallOption) (*DeviceSchemaIndexResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceSchemaIndex(ctx, in, opts...)
}

// 获取设备物模型列表
func (d *directDeviceManage) DeviceSchemaIndex(ctx context.Context, in *DeviceSchemaIndexReq, opts ...grpc.CallOption) (*DeviceSchemaIndexResp, error) {
	return d.svr.DeviceSchemaIndex(ctx, in)
}

func (m *defaultDeviceManage) DeviceSchemaTslRead(ctx context.Context, in *DeviceSchemaTslReadReq, opts ...grpc.CallOption) (*DeviceSchemaTslReadResp, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceSchemaTslRead(ctx, in, opts...)
}

func (d *directDeviceManage) DeviceSchemaTslRead(ctx context.Context, in *DeviceSchemaTslReadReq, opts ...grpc.CallOption) (*DeviceSchemaTslReadResp, error) {
	return d.svr.DeviceSchemaTslRead(ctx, in)
}

// 将设备加到多个分组中
func (m *defaultDeviceManage) DeviceGroupMultiCreate(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceGroupMultiCreate(ctx, in, opts...)
}

// 将设备加到多个分组中
func (d *directDeviceManage) DeviceGroupMultiCreate(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceGroupMultiCreate(ctx, in)
}

// 更新设备所在分组
func (m *defaultDeviceManage) DeviceGroupMultiUpdate(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceGroupMultiUpdate(ctx, in, opts...)
}

// 更新设备所在分组
func (d *directDeviceManage) DeviceGroupMultiUpdate(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceGroupMultiUpdate(ctx, in)
}

// 删除设备所在分组
func (m *defaultDeviceManage) DeviceGroupMultiDelete(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	client := dm.NewDeviceManageClient(m.cli.Conn())
	return client.DeviceGroupMultiDelete(ctx, in, opts...)
}

// 删除设备所在分组
func (d *directDeviceManage) DeviceGroupMultiDelete(ctx context.Context, in *DeviceGroupMultiSaveReq, opts ...grpc.CallOption) (*Empty, error) {
	return d.svr.DeviceGroupMultiDelete(ctx, in)
}
