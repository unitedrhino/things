package remoteconfiglogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoteConfigPushAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoteConfigPushAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoteConfigPushAllLogic {
	return &RemoteConfigPushAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoteConfigPushAllLogic) RemoteConfigPushAll(in *dm.RemoteConfigPushAllReq) (*dm.Response, error) {
	// 获取最后一条配置，并发布到设备

	return &dm.Response{}, nil
}
