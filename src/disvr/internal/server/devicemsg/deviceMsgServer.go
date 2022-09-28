// Code generated by goctl. DO NOT EDIT!
// Source: di.proto

package server

import (
	"context"

	"github.com/i-Things/things/src/disvr/internal/logic/devicemsg"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"
)

type DeviceMsgServer struct {
	svcCtx *svc.ServiceContext
	di.UnimplementedDeviceMsgServer
}

func NewDeviceMsgServer(svcCtx *svc.ServiceContext) *DeviceMsgServer {
	return &DeviceMsgServer{
		svcCtx: svcCtx,
	}
}

// 获取设备sdk调试日志
func (s *DeviceMsgServer) SdkLogIndex(ctx context.Context, in *di.SdkLogIndexReq) (*di.SdkLogIndexResp, error) {
	l := devicemsglogic.NewSdkLogIndexLogic(ctx, s.svcCtx)
	return l.SdkLogIndex(in)
}

// 获取设备调试信息记录登入登出,操作
func (s *DeviceMsgServer) HubLogIndex(ctx context.Context, in *di.HubLogIndexReq) (*di.HubLogIndexResp, error) {
	l := devicemsglogic.NewHubLogIndexLogic(ctx, s.svcCtx)
	return l.HubLogIndex(in)
}

// 获取设备数据信息
func (s *DeviceMsgServer) PropertyLatestIndex(ctx context.Context, in *di.PropertyLatestIndexReq) (*di.PropertyIndexResp, error) {
	l := devicemsglogic.NewPropertyLatestIndexLogic(ctx, s.svcCtx)
	return l.PropertyLatestIndex(in)
}

// 获取设备数据信息
func (s *DeviceMsgServer) PropertyLogIndex(ctx context.Context, in *di.PropertyLogIndexReq) (*di.PropertyIndexResp, error) {
	l := devicemsglogic.NewPropertyLogIndexLogic(ctx, s.svcCtx)
	return l.PropertyLogIndex(in)
}

// 获取设备数据信息
func (s *DeviceMsgServer) EventLogIndex(ctx context.Context, in *di.EventLogIndexReq) (*di.EventIndexResp, error) {
	l := devicemsglogic.NewEventLogIndexLogic(ctx, s.svcCtx)
	return l.EventLogIndex(in)
}
