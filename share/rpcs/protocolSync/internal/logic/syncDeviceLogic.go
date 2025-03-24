package logic

import (
	"context"

	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/internal/svc"
	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/pb/protocolSync"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDeviceLogic {
	return &SyncDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncDeviceLogic) SyncDevice(in *protocolSync.SyncDeviceReq) (*protocolSync.SyncDeviceResp, error) {
	// todo: add your logic here and delete this line

	return &protocolSync.SyncDeviceResp{}, nil
}
