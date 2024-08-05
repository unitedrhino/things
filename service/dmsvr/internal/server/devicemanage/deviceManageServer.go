// Code generated by goctl. DO NOT EDIT.
// Source: dm.proto

package server

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/logic/devicemanage"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
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
