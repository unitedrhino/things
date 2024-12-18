// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: dm.proto

package server

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic/userdevice"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

type UserDeviceServer struct {
	svcCtx *svc.ServiceContext
	dm.UnimplementedUserDeviceServer
}

func NewUserDeviceServer(svcCtx *svc.ServiceContext) *UserDeviceServer {
	return &UserDeviceServer{
		svcCtx: svcCtx,
	}
}

// 用户收藏的设备
func (s *UserDeviceServer) UserDeviceCollectMultiCreate(ctx context.Context, in *dm.UserDeviceCollectSave) (*dm.Empty, error) {
	l := userdevicelogic.NewUserDeviceCollectMultiCreateLogic(ctx, s.svcCtx)
	return l.UserDeviceCollectMultiCreate(in)
}

func (s *UserDeviceServer) UserDeviceCollectMultiDelete(ctx context.Context, in *dm.UserDeviceCollectSave) (*dm.Empty, error) {
	l := userdevicelogic.NewUserDeviceCollectMultiDeleteLogic(ctx, s.svcCtx)
	return l.UserDeviceCollectMultiDelete(in)
}

func (s *UserDeviceServer) UserDeviceCollectIndex(ctx context.Context, in *dm.Empty) (*dm.UserDeviceCollectSave, error) {
	l := userdevicelogic.NewUserDeviceCollectIndexLogic(ctx, s.svcCtx)
	return l.UserDeviceCollectIndex(in)
}

// 分享设备
func (s *UserDeviceServer) UserDeviceShareCreate(ctx context.Context, in *dm.UserDeviceShareInfo) (*dm.WithID, error) {
	l := userdevicelogic.NewUserDeviceShareCreateLogic(ctx, s.svcCtx)
	return l.UserDeviceShareCreate(in)
}

// 更新权限
func (s *UserDeviceServer) UserDeviceShareUpdate(ctx context.Context, in *dm.UserDeviceShareInfo) (*dm.Empty, error) {
	l := userdevicelogic.NewUserDeviceShareUpdateLogic(ctx, s.svcCtx)
	return l.UserDeviceShareUpdate(in)
}

// 取消分享设备
func (s *UserDeviceServer) UserDeviceShareDelete(ctx context.Context, in *dm.UserDeviceShareReadReq) (*dm.Empty, error) {
	l := userdevicelogic.NewUserDeviceShareDeleteLogic(ctx, s.svcCtx)
	return l.UserDeviceShareDelete(in)
}

// 取消分享设备
func (s *UserDeviceServer) UserDeviceShareMultiDelete(ctx context.Context, in *dm.UserDeviceShareMultiDeleteReq) (*dm.Empty, error) {
	l := userdevicelogic.NewUserDeviceShareMultiDeleteLogic(ctx, s.svcCtx)
	return l.UserDeviceShareMultiDelete(in)
}

// 获取设备分享列表(只有设备的所有者才能获取)
func (s *UserDeviceServer) UserDeviceShareIndex(ctx context.Context, in *dm.UserDeviceShareIndexReq) (*dm.UserDeviceShareIndexResp, error) {
	l := userdevicelogic.NewUserDeviceShareIndexLogic(ctx, s.svcCtx)
	return l.UserDeviceShareIndex(in)
}

// 获取设备分享的详情
func (s *UserDeviceServer) UserDeviceShareRead(ctx context.Context, in *dm.UserDeviceShareReadReq) (*dm.UserDeviceShareInfo, error) {
	l := userdevicelogic.NewUserDeviceShareReadLogic(ctx, s.svcCtx)
	return l.UserDeviceShareRead(in)
}

// 转让设备
func (s *UserDeviceServer) UserDeviceTransfer(ctx context.Context, in *dm.DeviceTransferReq) (*dm.Empty, error) {
	l := userdevicelogic.NewUserDeviceTransferLogic(ctx, s.svcCtx)
	return l.UserDeviceTransfer(in)
}

// rpc userDeviceOtaGetVersion(UserDeviceOtaGetVersionReq)returns(userDeviceOtaGetVersionResp);
func (s *UserDeviceServer) UserDeviceShareMultiCreate(ctx context.Context, in *dm.UserDeviceShareMultiInfo) (*dm.UserDeviceShareMultiToken, error) {
	l := userdevicelogic.NewUserDeviceShareMultiCreateLogic(ctx, s.svcCtx)
	return l.UserDeviceShareMultiCreate(in)
}

// 扫码后获取设备列表
func (s *UserDeviceServer) UserDeivceShareMultiIndex(ctx context.Context, in *dm.UserDeviceShareMultiToken) (*dm.UserDeviceShareMultiInfo, error) {
	l := userdevicelogic.NewUserDeivceShareMultiIndexLogic(ctx, s.svcCtx)
	return l.UserDeivceShareMultiIndex(in)
}

// 接受批量分享的设备
func (s *UserDeviceServer) UserDeivceShareMultiAccept(ctx context.Context, in *dm.UserDeviceShareMultiAcceptReq) (*dm.Empty, error) {
	l := userdevicelogic.NewUserDeivceShareMultiAcceptLogic(ctx, s.svcCtx)
	return l.UserDeivceShareMultiAccept(in)
}
