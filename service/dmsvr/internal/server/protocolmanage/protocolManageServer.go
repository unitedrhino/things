// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.1
// Source: dm.proto

package server

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic/protocolmanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
)

type ProtocolManageServer struct {
	svcCtx *svc.ServiceContext
	dm.UnimplementedProtocolManageServer
}

func NewProtocolManageServer(svcCtx *svc.ServiceContext) *ProtocolManageServer {
	return &ProtocolManageServer{
		svcCtx: svcCtx,
	}
}

// 协议列表
func (s *ProtocolManageServer) ProtocolInfoIndex(ctx context.Context, in *dm.ProtocolInfoIndexReq) (*dm.ProtocolInfoIndexResp, error) {
	l := protocolmanagelogic.NewProtocolInfoIndexLogic(ctx, s.svcCtx)
	return l.ProtocolInfoIndex(in)
}

// 协议详情
func (s *ProtocolManageServer) ProtocolInfoRead(ctx context.Context, in *dm.WithIDCode) (*dm.ProtocolInfo, error) {
	l := protocolmanagelogic.NewProtocolInfoReadLogic(ctx, s.svcCtx)
	return l.ProtocolInfoRead(in)
}

// 协议创建
func (s *ProtocolManageServer) ProtocolInfoCreate(ctx context.Context, in *dm.ProtocolInfo) (*dm.WithID, error) {
	l := protocolmanagelogic.NewProtocolInfoCreateLogic(ctx, s.svcCtx)
	return l.ProtocolInfoCreate(in)
}

// 协议更新
func (s *ProtocolManageServer) ProtocolInfoUpdate(ctx context.Context, in *dm.ProtocolInfo) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolInfoUpdateLogic(ctx, s.svcCtx)
	return l.ProtocolInfoUpdate(in)
}

// 协议删除
func (s *ProtocolManageServer) ProtocolInfoDelete(ctx context.Context, in *dm.WithID) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolInfoDeleteLogic(ctx, s.svcCtx)
	return l.ProtocolInfoDelete(in)
}

// 更新服务状态,只给服务调用
func (s *ProtocolManageServer) ProtocolServiceUpdate(ctx context.Context, in *dm.ProtocolService) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolServiceUpdateLogic(ctx, s.svcCtx)
	return l.ProtocolServiceUpdate(in)
}

func (s *ProtocolManageServer) ProtocolServiceDelete(ctx context.Context, in *dm.WithID) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolServiceDeleteLogic(ctx, s.svcCtx)
	return l.ProtocolServiceDelete(in)
}

func (s *ProtocolManageServer) ProtocolServiceIndex(ctx context.Context, in *dm.ProtocolServiceIndexReq) (*dm.ProtocolServiceIndexResp, error) {
	l := protocolmanagelogic.NewProtocolServiceIndexLogic(ctx, s.svcCtx)
	return l.ProtocolServiceIndex(in)
}

// 协议列表
func (s *ProtocolManageServer) ProtocolScriptIndex(ctx context.Context, in *dm.ProtocolScriptIndexReq) (*dm.ProtocolScriptIndexResp, error) {
	l := protocolmanagelogic.NewProtocolScriptIndexLogic(ctx, s.svcCtx)
	return l.ProtocolScriptIndex(in)
}

// 协议详情
func (s *ProtocolManageServer) ProtocolScriptRead(ctx context.Context, in *dm.WithID) (*dm.ProtocolScript, error) {
	l := protocolmanagelogic.NewProtocolScriptReadLogic(ctx, s.svcCtx)
	return l.ProtocolScriptRead(in)
}

// 协议创建
func (s *ProtocolManageServer) ProtocolScriptCreate(ctx context.Context, in *dm.ProtocolScript) (*dm.WithID, error) {
	l := protocolmanagelogic.NewProtocolScriptCreateLogic(ctx, s.svcCtx)
	return l.ProtocolScriptCreate(in)
}

// 协议更新
func (s *ProtocolManageServer) ProtocolScriptUpdate(ctx context.Context, in *dm.ProtocolScript) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolScriptUpdateLogic(ctx, s.svcCtx)
	return l.ProtocolScriptUpdate(in)
}

// 协议删除
func (s *ProtocolManageServer) ProtocolScriptDelete(ctx context.Context, in *dm.WithID) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolScriptDeleteLogic(ctx, s.svcCtx)
	return l.ProtocolScriptDelete(in)
}

func (s *ProtocolManageServer) ProtocolScriptDebug(ctx context.Context, in *dm.ProtocolScriptDebugReq) (*dm.ProtocolScriptDebugResp, error) {
	l := protocolmanagelogic.NewProtocolScriptDebugLogic(ctx, s.svcCtx)
	return l.ProtocolScriptDebug(in)
}

// 协议列表
func (s *ProtocolManageServer) ProtocolScriptDeviceIndex(ctx context.Context, in *dm.ProtocolScriptDeviceIndexReq) (*dm.ProtocolScriptDeviceIndexResp, error) {
	l := protocolmanagelogic.NewProtocolScriptDeviceIndexLogic(ctx, s.svcCtx)
	return l.ProtocolScriptDeviceIndex(in)
}

// 协议详情
func (s *ProtocolManageServer) ProtocolScriptDeviceRead(ctx context.Context, in *dm.WithID) (*dm.ProtocolScriptDevice, error) {
	l := protocolmanagelogic.NewProtocolScriptDeviceReadLogic(ctx, s.svcCtx)
	return l.ProtocolScriptDeviceRead(in)
}

// 协议创建
func (s *ProtocolManageServer) ProtocolScriptDeviceCreate(ctx context.Context, in *dm.ProtocolScriptDevice) (*dm.WithID, error) {
	l := protocolmanagelogic.NewProtocolScriptDeviceCreateLogic(ctx, s.svcCtx)
	return l.ProtocolScriptDeviceCreate(in)
}

// 协议更新
func (s *ProtocolManageServer) ProtocolScriptDeviceUpdate(ctx context.Context, in *dm.ProtocolScriptDevice) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolScriptDeviceUpdateLogic(ctx, s.svcCtx)
	return l.ProtocolScriptDeviceUpdate(in)
}

// 协议删除
func (s *ProtocolManageServer) ProtocolScriptDeviceDelete(ctx context.Context, in *dm.WithID) (*dm.Empty, error) {
	l := protocolmanagelogic.NewProtocolScriptDeviceDeleteLogic(ctx, s.svcCtx)
	return l.ProtocolScriptDeviceDelete(in)
}

func (s *ProtocolManageServer) ProtocolScriptMultiImport(ctx context.Context, in *dm.ProtocolScriptImportReq) (*dm.ImportResp, error) {
	l := protocolmanagelogic.NewProtocolScriptMultiImportLogic(ctx, s.svcCtx)
	return l.ProtocolScriptMultiImport(in)
}

func (s *ProtocolManageServer) ProtocolScriptMultiExport(ctx context.Context, in *dm.ProtocolScriptExportReq) (*dm.ProtocolScriptExportResp, error) {
	l := protocolmanagelogic.NewProtocolScriptMultiExportLogic(ctx, s.svcCtx)
	return l.ProtocolScriptMultiExport(in)
}
