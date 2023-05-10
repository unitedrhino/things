// Code generated by goctl. DO NOT EDIT.
// Source: dm.proto

package client

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AccessAuthReq               = dm.AccessAuthReq
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
	DeviceTypeCountReq          = dm.DeviceTypeCountReq
	DeviceTypeCountResp         = dm.DeviceTypeCountResp
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
	LoginAuthReq                = dm.LoginAuthReq
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
	RemoteConfigCreateReq       = dm.RemoteConfigCreateReq
	RemoteConfigIndexReq        = dm.RemoteConfigIndexReq
	RemoteConfigIndexResp       = dm.RemoteConfigIndexResp
	RemoteConfigLastReadReq     = dm.RemoteConfigLastReadReq
	RemoteConfigLastReadResp    = dm.RemoteConfigLastReadResp
	RemoteConfigPushAllReq      = dm.RemoteConfigPushAllReq
	Response                    = dm.Response
	RootCheckReq                = dm.RootCheckReq

	RemoteConfig interface {
		RemoteConfigCreate(ctx context.Context, in *RemoteConfigCreateReq, opts ...grpc.CallOption) (*Response, error)
		RemoteConfigIndex(ctx context.Context, in *RemoteConfigIndexReq, opts ...grpc.CallOption) (*RemoteConfigIndexResp, error)
		RemoteConfigPushAll(ctx context.Context, in *RemoteConfigPushAllReq, opts ...grpc.CallOption) (*Response, error)
		RemoteConfigLastRead(ctx context.Context, in *RemoteConfigLastReadReq, opts ...grpc.CallOption) (*RemoteConfigLastReadResp, error)
	}

	defaultRemoteConfig struct {
		cli zrpc.Client
	}

	directRemoteConfig struct {
		svcCtx *svc.ServiceContext
		svr    dm.RemoteConfigServer
	}
)

func NewRemoteConfig(cli zrpc.Client) RemoteConfig {
	return &defaultRemoteConfig{
		cli: cli,
	}
}

func NewDirectRemoteConfig(svcCtx *svc.ServiceContext, svr dm.RemoteConfigServer) RemoteConfig {
	return &directRemoteConfig{
		svr:    svr,
		svcCtx: svcCtx,
	}
}

func (m *defaultRemoteConfig) RemoteConfigCreate(ctx context.Context, in *RemoteConfigCreateReq, opts ...grpc.CallOption) (*Response, error) {
	client := dm.NewRemoteConfigClient(m.cli.Conn())
	return client.RemoteConfigCreate(ctx, in, opts...)
}

func (d *directRemoteConfig) RemoteConfigCreate(ctx context.Context, in *RemoteConfigCreateReq, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.RemoteConfigCreate(ctx, in)
}

func (m *defaultRemoteConfig) RemoteConfigIndex(ctx context.Context, in *RemoteConfigIndexReq, opts ...grpc.CallOption) (*RemoteConfigIndexResp, error) {
	client := dm.NewRemoteConfigClient(m.cli.Conn())
	return client.RemoteConfigIndex(ctx, in, opts...)
}

func (d *directRemoteConfig) RemoteConfigIndex(ctx context.Context, in *RemoteConfigIndexReq, opts ...grpc.CallOption) (*RemoteConfigIndexResp, error) {
	return d.svr.RemoteConfigIndex(ctx, in)
}

func (m *defaultRemoteConfig) RemoteConfigPushAll(ctx context.Context, in *RemoteConfigPushAllReq, opts ...grpc.CallOption) (*Response, error) {
	client := dm.NewRemoteConfigClient(m.cli.Conn())
	return client.RemoteConfigPushAll(ctx, in, opts...)
}

func (d *directRemoteConfig) RemoteConfigPushAll(ctx context.Context, in *RemoteConfigPushAllReq, opts ...grpc.CallOption) (*Response, error) {
	return d.svr.RemoteConfigPushAll(ctx, in)
}

func (m *defaultRemoteConfig) RemoteConfigLastRead(ctx context.Context, in *RemoteConfigLastReadReq, opts ...grpc.CallOption) (*RemoteConfigLastReadResp, error) {
	client := dm.NewRemoteConfigClient(m.cli.Conn())
	return client.RemoteConfigLastRead(ctx, in, opts...)
}

func (d *directRemoteConfig) RemoteConfigLastRead(ctx context.Context, in *RemoteConfigLastReadReq, opts ...grpc.CallOption) (*RemoteConfigLastReadResp, error) {
	return d.svr.RemoteConfigLastRead(ctx, in)
}
