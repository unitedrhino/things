package interact

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"sync"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyControlMultiSendLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	retMsg []*types.DeviceInteractMultiSendPropertyMsg
	err    error
	mutex  sync.Mutex
}

func NewPropertyControlMultiSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyControlMultiSendLogic {
	return &PropertyControlMultiSendLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PropertyControlMultiSendLogic) PropertyControlMultiSend(req *types.DeviceInteractMultiSendPropertyReq) (resp *types.DeviceInteractMultiSendPropertyResp, err error) {
	ret, err := l.svcCtx.DeviceInteract.PropertyControlMultiSend(l.ctx, utils.Copy[dm.PropertyControlMultiSendReq](req))
	if err != nil {
		return nil, err
	}
	return utils.Copy[types.DeviceInteractMultiSendPropertyResp](ret), nil
}
