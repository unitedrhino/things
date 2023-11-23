// Code generated by goctl. DO NOT EDIT.
// Source: dm.proto

package devicegroup

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AccessAuthReq               = dm.AccessAuthReq
	DeleteOtaFirmwareReq        = dm.DeleteOtaFirmwareReq
	DeviceCore                  = dm.DeviceCore
	DeviceGatewayBindDevice     = dm.DeviceGatewayBindDevice
	DeviceGatewayIndexReq       = dm.DeviceGatewayIndexReq
	DeviceGatewayIndexResp      = dm.DeviceGatewayIndexResp
	DeviceGatewayMultiCreateReq = dm.DeviceGatewayMultiCreateReq
	DeviceGatewayMultiDeleteReq = dm.DeviceGatewayMultiDeleteReq
	DeviceGatewaySign           = dm.DeviceGatewaySign
	DeviceInfo                  = dm.DeviceInfo
	DeviceInfoCountReq          = dm.DeviceInfoCountReq
	DeviceInfoCountResp         = dm.DeviceInfoCountResp
	DeviceInfoDeleteReq         = dm.DeviceInfoDeleteReq
	DeviceInfoIndexReq          = dm.DeviceInfoIndexReq
	DeviceInfoIndexResp         = dm.DeviceInfoIndexResp
	DeviceInfoReadReq           = dm.DeviceInfoReadReq
	DeviceRegisterReq           = dm.DeviceRegisterReq
	DeviceRegisterResp          = dm.DeviceRegisterResp
	DeviceTypeCountReq          = dm.DeviceTypeCountReq
	DeviceTypeCountResp         = dm.DeviceTypeCountResp
	EventIndex                  = dm.EventIndex
	EventIndexResp              = dm.EventIndexResp
	EventLogIndexReq            = dm.EventLogIndexReq
	Firmware                    = dm.Firmware
	FirmwareFile                = dm.FirmwareFile
	FirmwareInfo                = dm.FirmwareInfo
	FirmwareInfoDeleteReq       = dm.FirmwareInfoDeleteReq
	FirmwareInfoDeleteResp      = dm.FirmwareInfoDeleteResp
	FirmwareInfoIndexReq        = dm.FirmwareInfoIndexReq
	FirmwareInfoIndexResp       = dm.FirmwareInfoIndexResp
	FirmwareInfoReadReq         = dm.FirmwareInfoReadReq
	FirmwareInfoReadResp        = dm.FirmwareInfoReadResp
	FirmwareResp                = dm.FirmwareResp
	GetPropertyReplyReq         = dm.GetPropertyReplyReq
	GetPropertyReplyResp        = dm.GetPropertyReplyResp
	GroupDeviceIndexReq         = dm.GroupDeviceIndexReq
	GroupDeviceIndexResp        = dm.GroupDeviceIndexResp
	GroupDeviceMultiCreateReq   = dm.GroupDeviceMultiCreateReq
	GroupDeviceMultiDeleteReq   = dm.GroupDeviceMultiDeleteReq
	GroupInfo                   = dm.GroupInfo
	GroupInfoCreateReq          = dm.GroupInfoCreateReq
	GroupInfoDeleteReq          = dm.GroupInfoDeleteReq
	GroupInfoIndexReq           = dm.GroupInfoIndexReq
	GroupInfoIndexResp          = dm.GroupInfoIndexResp
	GroupInfoReadReq            = dm.GroupInfoReadReq
	GroupInfoUpdateReq          = dm.GroupInfoUpdateReq
	HubLogIndex                 = dm.HubLogIndex
	HubLogIndexReq              = dm.HubLogIndexReq
	HubLogIndexResp             = dm.HubLogIndexResp
	ListOtaFirmwareReq          = dm.ListOtaFirmwareReq
	ListOtaFirmwareResp         = dm.ListOtaFirmwareResp
	LoginAuthReq                = dm.LoginAuthReq
	ModifyOtaFirmwareReq        = dm.ModifyOtaFirmwareReq
	MultiSendPropertyReq        = dm.MultiSendPropertyReq
	MultiSendPropertyResp       = dm.MultiSendPropertyResp
	OtaCommonResp               = dm.OtaCommonResp
	OtaFirmwareDeviceInfoReq    = dm.OtaFirmwareDeviceInfoReq
	OtaFirmwareDeviceInfoResp   = dm.OtaFirmwareDeviceInfoResp
	OtaFirmwareFile             = dm.OtaFirmwareFile
	OtaFirmwareFileIndexReq     = dm.OtaFirmwareFileIndexReq
	OtaFirmwareFileIndexResp    = dm.OtaFirmwareFileIndexResp
	OtaFirmwareFileInfo         = dm.OtaFirmwareFileInfo
	OtaFirmwareFileReq          = dm.OtaFirmwareFileReq
	OtaFirmwareFileResp         = dm.OtaFirmwareFileResp
	OtaFirmwareInfo             = dm.OtaFirmwareInfo
	OtaFirmwareReq              = dm.OtaFirmwareReq
	OtaFirmwareResp             = dm.OtaFirmwareResp
	OtaPageInfo                 = dm.OtaPageInfo
	OtaPromptIndexReq           = dm.OtaPromptIndexReq
	OtaPromptIndexResp          = dm.OtaPromptIndexResp
	OtaTaskBatchReq             = dm.OtaTaskBatchReq
	OtaTaskBatchResp            = dm.OtaTaskBatchResp
	OtaTaskCancleReq            = dm.OtaTaskCancleReq
	OtaTaskCreatResp            = dm.OtaTaskCreatResp
	OtaTaskCreateReq            = dm.OtaTaskCreateReq
	OtaTaskDeviceCancleReq      = dm.OtaTaskDeviceCancleReq
	OtaTaskDeviceIndexReq       = dm.OtaTaskDeviceIndexReq
	OtaTaskDeviceIndexResp      = dm.OtaTaskDeviceIndexResp
	OtaTaskDeviceInfo           = dm.OtaTaskDeviceInfo
	OtaTaskDeviceProcessReq     = dm.OtaTaskDeviceProcessReq
	OtaTaskDeviceReadReq        = dm.OtaTaskDeviceReadReq
	OtaTaskIndexReq             = dm.OtaTaskIndexReq
	OtaTaskIndexResp            = dm.OtaTaskIndexResp
	OtaTaskInfo                 = dm.OtaTaskInfo
	OtaTaskReadReq              = dm.OtaTaskReadReq
	OtaTaskReadResp             = dm.OtaTaskReadResp
	PageInfo                    = dm.PageInfo
	PageInfo_OrderBy            = dm.PageInfo_OrderBy
	Point                       = dm.Point
	ProductCustom               = dm.ProductCustom
	ProductCustomReadReq        = dm.ProductCustomReadReq
	ProductInfo                 = dm.ProductInfo
	ProductInfoDeleteReq        = dm.ProductInfoDeleteReq
	ProductInfoIndexReq         = dm.ProductInfoIndexReq
	ProductInfoIndexResp        = dm.ProductInfoIndexResp
	ProductInfoReadReq          = dm.ProductInfoReadReq
	ProductRemoteConfig         = dm.ProductRemoteConfig
	ProductSchemaCreateReq      = dm.ProductSchemaCreateReq
	ProductSchemaDeleteReq      = dm.ProductSchemaDeleteReq
	ProductSchemaIndexReq       = dm.ProductSchemaIndexReq
	ProductSchemaIndexResp      = dm.ProductSchemaIndexResp
	ProductSchemaInfo           = dm.ProductSchemaInfo
	ProductSchemaTslImportReq   = dm.ProductSchemaTslImportReq
	ProductSchemaTslReadReq     = dm.ProductSchemaTslReadReq
	ProductSchemaTslReadResp    = dm.ProductSchemaTslReadResp
	ProductSchemaUpdateReq      = dm.ProductSchemaUpdateReq
	PropertyIndex               = dm.PropertyIndex
	PropertyIndexResp           = dm.PropertyIndexResp
	PropertyLatestIndexReq      = dm.PropertyLatestIndexReq
	PropertyLogIndexReq         = dm.PropertyLogIndexReq
	QueryOtaFirmwareReq         = dm.QueryOtaFirmwareReq
	QueryOtaFirmwareResp        = dm.QueryOtaFirmwareResp
	RemoteConfigCreateReq       = dm.RemoteConfigCreateReq
	RemoteConfigIndexReq        = dm.RemoteConfigIndexReq
	RemoteConfigIndexResp       = dm.RemoteConfigIndexResp
	RemoteConfigLastReadReq     = dm.RemoteConfigLastReadReq
	RemoteConfigLastReadResp    = dm.RemoteConfigLastReadResp
	RemoteConfigPushAllReq      = dm.RemoteConfigPushAllReq
	RespActionReq               = dm.RespActionReq
	RespReadReq                 = dm.RespReadReq
	Response                    = dm.Response
	RootCheckReq                = dm.RootCheckReq
	SdkLogIndex                 = dm.SdkLogIndex
	SdkLogIndexReq              = dm.SdkLogIndexReq
	SdkLogIndexResp             = dm.SdkLogIndexResp
	SendActionReq               = dm.SendActionReq
	SendActionResp              = dm.SendActionResp
	SendMsgReq                  = dm.SendMsgReq
	SendMsgResp                 = dm.SendMsgResp
	SendOption                  = dm.SendOption
	SendPropertyMsg             = dm.SendPropertyMsg
	SendPropertyReq             = dm.SendPropertyReq
	SendPropertyResp            = dm.SendPropertyResp
	ShadowIndex                 = dm.ShadowIndex
	ShadowIndexResp             = dm.ShadowIndexResp
	Tag                         = dm.Tag
	VerifyOtaFirmwareReq        = dm.VerifyOtaFirmwareReq

	DeviceGroup interface {
		// 创建分组
		GroupInfoCreate(ctx context.Context, in *GroupInfoCreateReq, opts ...grpc.CallOption) (*Response, error)
		// 获取分组信息列表
		GroupInfoIndex(ctx context.Context, in *GroupInfoIndexReq, opts ...grpc.CallOption) (*GroupInfoIndexResp, error)
		// 获取分组信息详情
		GroupInfoRead(ctx context.Context, in *GroupInfoReadReq, opts ...grpc.CallOption) (*GroupInfo, error)
		// 更新分组
		GroupInfoUpdate(ctx context.Context, in *GroupInfoUpdateReq, opts ...grpc.CallOption) (*Response, error)
		// 删除分组
		GroupInfoDelete(ctx context.Context, in *GroupInfoDeleteReq, opts ...grpc.CallOption) (*Response, error)
		// 创建分组设备
		GroupDeviceMultiCreate(ctx context.Context, in *GroupDeviceMultiCreateReq, opts ...grpc.CallOption) (*Response, error)
		// 获取分组设备信息列表
		GroupDeviceIndex(ctx context.Context, in *GroupDeviceIndexReq, opts ...grpc.CallOption) (*GroupDeviceIndexResp, error)
		// 删除分组设备
		GroupDeviceMultiDelete(ctx context.Context, in *GroupDeviceMultiDeleteReq, opts ...grpc.CallOption) (*Response, error)
	}

	defaultDeviceGroup struct {
		cli zrpc.Client
	}

	directDeviceGroup struct {
		svcCtx *svc.ServiceContext
		svr    dm.DeviceGroupServer
	}
)

func NewDeviceGroup(cli zrpc.Client) DeviceGroup {
	return &defaultDeviceGroup{
		cli: cli,
	}
}

func NewDirectDeviceGroup(svcCtx *svc.ServiceContext, svr dm.DeviceGroupServer) DeviceGroup {
	return &directDeviceGroup{
		svr:    svr,
		svcCtx: svcCtx,
	}
}

// 创建分组
func (m *defaultDeviceGroup) GroupInfoCreate(ctx context.Context, in *GroupInfoCreateReq, opts ...grpc.CallOption) (*Response, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupInfoCreate(ctx, in, opts...)
}

// 创建分组
func (d *directDeviceGroup) GroupInfoCreate(ctx context.Context, in *GroupInfoCreateReq, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.GroupInfoCreate(ctx, in)
}

// 获取分组信息列表
func (m *defaultDeviceGroup) GroupInfoIndex(ctx context.Context, in *GroupInfoIndexReq, opts ...grpc.CallOption) (*GroupInfoIndexResp, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupInfoIndex(ctx, in, opts...)
}

// 获取分组信息列表
func (d *directDeviceGroup) GroupInfoIndex(ctx context.Context, in *GroupInfoIndexReq, opts ...grpc.CallOption) (*GroupInfoIndexResp, error) {
	return d.svr.GroupInfoIndex(ctx, in)
}

// 获取分组信息详情
func (m *defaultDeviceGroup) GroupInfoRead(ctx context.Context, in *GroupInfoReadReq, opts ...grpc.CallOption) (*GroupInfo, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupInfoRead(ctx, in, opts...)
}

// 获取分组信息详情
func (d *directDeviceGroup) GroupInfoRead(ctx context.Context, in *GroupInfoReadReq, opts ...grpc.CallOption) (*GroupInfo, error) {
	return d.svr.GroupInfoRead(ctx, in)
}

// 更新分组
func (m *defaultDeviceGroup) GroupInfoUpdate(ctx context.Context, in *GroupInfoUpdateReq, opts ...grpc.CallOption) (*Response, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupInfoUpdate(ctx, in, opts...)
}

// 更新分组
func (d *directDeviceGroup) GroupInfoUpdate(ctx context.Context, in *GroupInfoUpdateReq, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.GroupInfoUpdate(ctx, in)
}

// 删除分组
func (m *defaultDeviceGroup) GroupInfoDelete(ctx context.Context, in *GroupInfoDeleteReq, opts ...grpc.CallOption) (*Response, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupInfoDelete(ctx, in, opts...)
}

// 删除分组
func (d *directDeviceGroup) GroupInfoDelete(ctx context.Context, in *GroupInfoDeleteReq, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.GroupInfoDelete(ctx, in)
}

// 创建分组设备
func (m *defaultDeviceGroup) GroupDeviceMultiCreate(ctx context.Context, in *GroupDeviceMultiCreateReq, opts ...grpc.CallOption) (*Response, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupDeviceMultiCreate(ctx, in, opts...)
}

// 创建分组设备
func (d *directDeviceGroup) GroupDeviceMultiCreate(ctx context.Context, in *GroupDeviceMultiCreateReq, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.GroupDeviceMultiCreate(ctx, in)
}

// 获取分组设备信息列表
func (m *defaultDeviceGroup) GroupDeviceIndex(ctx context.Context, in *GroupDeviceIndexReq, opts ...grpc.CallOption) (*GroupDeviceIndexResp, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupDeviceIndex(ctx, in, opts...)
}

// 获取分组设备信息列表
func (d *directDeviceGroup) GroupDeviceIndex(ctx context.Context, in *GroupDeviceIndexReq, opts ...grpc.CallOption) (*GroupDeviceIndexResp, error) {
	return d.svr.GroupDeviceIndex(ctx, in)
}

// 删除分组设备
func (m *defaultDeviceGroup) GroupDeviceMultiDelete(ctx context.Context, in *GroupDeviceMultiDeleteReq, opts ...grpc.CallOption) (*Response, error) {
	client := dm.NewDeviceGroupClient(m.cli.Conn())
	return client.GroupDeviceMultiDelete(ctx, in, opts...)
}

// 删除分组设备
func (d *directDeviceGroup) GroupDeviceMultiDelete(ctx context.Context, in *GroupDeviceMultiDeleteReq, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.GroupDeviceMultiDelete(ctx, in)
}
