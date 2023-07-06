package remoteconfiglogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
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
	err := l.svcCtx.DataUpdate.DeviceRemoteConfigUpdate(l.ctx, &events.DeviceUpdateInfo{
		ProductID: in.ProductID,
	})
	if err != nil {
		l.Errorf("RemoteConfigPushAll.DeviceRemoteConfigUpdate err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.Response{}, nil
}
