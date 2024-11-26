// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: dm.proto

package server

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

type DeviceManageServer struct {
	svcCtx *svc.ServiceContext
	dm.UnimplementedDeviceManageServer
}

func NewDeviceManageServer(svcCtx *svc.ServiceContext) *DeviceManageServer {
	return &DeviceManageServer{
		svcCtx: svcCtx,
	}
}

// 鉴定是否是root账号(提供给mqtt broker)
func (s *DeviceManageServer) RootCheck(ctx context.Context, in *dm.RootCheckReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewRootCheckLogic(ctx, s.svcCtx)
	return l.RootCheck(in)
}

// 新增设备
func (s *DeviceManageServer) DeviceInfoCreate(ctx context.Context, in *dm.DeviceInfo) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceInfoCreateLogic(ctx, s.svcCtx)
	return l.DeviceInfoCreate(in)
}

// 更新设备
func (s *DeviceManageServer) DeviceInfoUpdate(ctx context.Context, in *dm.DeviceInfo) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceInfoUpdateLogic(ctx, s.svcCtx)
	return l.DeviceInfoUpdate(in)
}

func (s *DeviceManageServer) DeviceOnlineMultiFix(ctx context.Context, in *dm.DeviceOnlineMultiFixReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceOnlineMultiFixLogic(ctx, s.svcCtx)
	return l.DeviceOnlineMultiFix(in)
}

// 删除设备
func (s *DeviceManageServer) DeviceInfoDelete(ctx context.Context, in *dm.DeviceInfoDeleteReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceInfoDeleteLogic(ctx, s.svcCtx)
	return l.DeviceInfoDelete(in)
}

// 获取设备信息列表
func (s *DeviceManageServer) DeviceInfoIndex(ctx context.Context, in *dm.DeviceInfoIndexReq) (*dm.DeviceInfoIndexResp, error) {
	l := devicemanagelogic.NewDeviceInfoIndexLogic(ctx, s.svcCtx)
	return l.DeviceInfoIndex(in)
}

// 批量更新设备状态
func (s *DeviceManageServer) DeviceInfoMultiUpdate(ctx context.Context, in *dm.DeviceInfoMultiUpdateReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceInfoMultiUpdateLogic(ctx, s.svcCtx)
	return l.DeviceInfoMultiUpdate(in)
}

// 获取设备信息详情
func (s *DeviceManageServer) DeviceInfoRead(ctx context.Context, in *dm.DeviceInfoReadReq) (*dm.DeviceInfo, error) {
	l := devicemanagelogic.NewDeviceInfoReadLogic(ctx, s.svcCtx)
	return l.DeviceInfoRead(in)
}

func (s *DeviceManageServer) DeviceInfoBind(ctx context.Context, in *dm.DeviceInfoBindReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceInfoBindLogic(ctx, s.svcCtx)
	return l.DeviceInfoBind(in)
}

func (s *DeviceManageServer) DeviceInfoMultiBind(ctx context.Context, in *dm.DeviceInfoMultiBindReq) (*dm.DeviceInfoMultiBindResp, error) {
	l := devicemanagelogic.NewDeviceInfoMultiBindLogic(ctx, s.svcCtx)
	return l.DeviceInfoMultiBind(in)
}

func (s *DeviceManageServer) DeviceInfoCanBind(ctx context.Context, in *dm.DeviceInfoCanBindReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceInfoCanBindLogic(ctx, s.svcCtx)
	return l.DeviceInfoCanBind(in)
}

func (s *DeviceManageServer) DeviceInfoUnbind(ctx context.Context, in *dm.DeviceCore) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceInfoUnbindLogic(ctx, s.svcCtx)
	return l.DeviceInfoUnbind(in)
}

func (s *DeviceManageServer) DeviceTransfer(ctx context.Context, in *dm.DeviceTransferReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceTransferLogic(ctx, s.svcCtx)
	return l.DeviceTransfer(in)
}

func (s *DeviceManageServer) DeviceReset(ctx context.Context, in *dm.DeviceResetReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceResetLogic(ctx, s.svcCtx)
	return l.DeviceReset(in)
}

func (s *DeviceManageServer) DeviceMove(ctx context.Context, in *dm.DeviceMoveReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceMoveLogic(ctx, s.svcCtx)
	return l.DeviceMove(in)
}

func (s *DeviceManageServer) DeviceModuleVersionRead(ctx context.Context, in *dm.DeviceModuleVersionReadReq) (*dm.DeviceModuleVersion, error) {
	l := devicemanagelogic.NewDeviceModuleVersionReadLogic(ctx, s.svcCtx)
	return l.DeviceModuleVersionRead(in)
}

func (s *DeviceManageServer) DeviceModuleVersionIndex(ctx context.Context, in *dm.DeviceModuleVersionIndexReq) (*dm.DeviceModuleVersionIndexResp, error) {
	l := devicemanagelogic.NewDeviceModuleVersionIndexLogic(ctx, s.svcCtx)
	return l.DeviceModuleVersionIndex(in)
}

// 绑定网关下子设备设备
func (s *DeviceManageServer) DeviceGatewayMultiCreate(ctx context.Context, in *dm.DeviceGatewayMultiCreateReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceGatewayMultiCreateLogic(ctx, s.svcCtx)
	return l.DeviceGatewayMultiCreate(in)
}

// 绑定网关下子设备设备
func (s *DeviceManageServer) DeviceGatewayMultiUpdate(ctx context.Context, in *dm.DeviceGatewayMultiSaveReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceGatewayMultiUpdateLogic(ctx, s.svcCtx)
	return l.DeviceGatewayMultiUpdate(in)
}

// 获取绑定信息的设备信息列表
func (s *DeviceManageServer) DeviceGatewayIndex(ctx context.Context, in *dm.DeviceGatewayIndexReq) (*dm.DeviceGatewayIndexResp, error) {
	l := devicemanagelogic.NewDeviceGatewayIndexLogic(ctx, s.svcCtx)
	return l.DeviceGatewayIndex(in)
}

// 删除网关下子设备
func (s *DeviceManageServer) DeviceGatewayMultiDelete(ctx context.Context, in *dm.DeviceGatewayMultiSaveReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceGatewayMultiDeleteLogic(ctx, s.svcCtx)
	return l.DeviceGatewayMultiDelete(in)
}

// 设备计数
func (s *DeviceManageServer) DeviceInfoCount(ctx context.Context, in *dm.DeviceInfoCountReq) (*dm.DeviceInfoCount, error) {
	l := devicemanagelogic.NewDeviceInfoCountLogic(ctx, s.svcCtx)
	return l.DeviceInfoCount(in)
}

// 设备类型
func (s *DeviceManageServer) DeviceTypeCount(ctx context.Context, in *dm.DeviceTypeCountReq) (*dm.DeviceTypeCountResp, error) {
	l := devicemanagelogic.NewDeviceTypeCountLogic(ctx, s.svcCtx)
	return l.DeviceTypeCount(in)
}

func (s *DeviceManageServer) DeviceCount(ctx context.Context, in *dm.DeviceCountReq) (*dm.DeviceCountResp, error) {
	l := devicemanagelogic.NewDeviceCountLogic(ctx, s.svcCtx)
	return l.DeviceCount(in)
}

func (s *DeviceManageServer) DeviceProfileRead(ctx context.Context, in *dm.DeviceProfileReadReq) (*dm.DeviceProfile, error) {
	l := devicemanagelogic.NewDeviceProfileReadLogic(ctx, s.svcCtx)
	return l.DeviceProfileRead(in)
}

func (s *DeviceManageServer) DeviceProfileDelete(ctx context.Context, in *dm.DeviceProfileReadReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceProfileDeleteLogic(ctx, s.svcCtx)
	return l.DeviceProfileDelete(in)
}

func (s *DeviceManageServer) DeviceProfileUpdate(ctx context.Context, in *dm.DeviceProfile) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceProfileUpdateLogic(ctx, s.svcCtx)
	return l.DeviceProfileUpdate(in)
}

func (s *DeviceManageServer) DeviceProfileIndex(ctx context.Context, in *dm.DeviceProfileIndexReq) (*dm.DeviceProfileIndexResp, error) {
	l := devicemanagelogic.NewDeviceProfileIndexLogic(ctx, s.svcCtx)
	return l.DeviceProfileIndex(in)
}

// 更新设备物模型
func (s *DeviceManageServer) DeviceSchemaUpdate(ctx context.Context, in *dm.DeviceSchema) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceSchemaUpdateLogic(ctx, s.svcCtx)
	return l.DeviceSchemaUpdate(in)
}

// 新增设备
func (s *DeviceManageServer) DeviceSchemaCreate(ctx context.Context, in *dm.DeviceSchema) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceSchemaCreateLogic(ctx, s.svcCtx)
	return l.DeviceSchemaCreate(in)
}

// 批量新增物模型,只新增没有的,已有的不处理
func (s *DeviceManageServer) DeviceSchemaMultiCreate(ctx context.Context, in *dm.DeviceSchemaMultiCreateReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceSchemaMultiCreateLogic(ctx, s.svcCtx)
	return l.DeviceSchemaMultiCreate(in)
}

// 删除设备物模型
func (s *DeviceManageServer) DeviceSchemaMultiDelete(ctx context.Context, in *dm.DeviceSchemaMultiDeleteReq) (*dm.Empty, error) {
	l := devicemanagelogic.NewDeviceSchemaMultiDeleteLogic(ctx, s.svcCtx)
	return l.DeviceSchemaMultiDelete(in)
}

// 获取设备物模型列表
func (s *DeviceManageServer) DeviceSchemaIndex(ctx context.Context, in *dm.DeviceSchemaIndexReq) (*dm.DeviceSchemaIndexResp, error) {
	l := devicemanagelogic.NewDeviceSchemaIndexLogic(ctx, s.svcCtx)
	return l.DeviceSchemaIndex(in)
}

func (s *DeviceManageServer) DeviceSchemaTslRead(ctx context.Context, in *dm.DeviceSchemaTslReadReq) (*dm.DeviceSchemaTslReadResp, error) {
	l := devicemanagelogic.NewDeviceSchemaTslReadLogic(ctx, s.svcCtx)
	return l.DeviceSchemaTslRead(in)
}
