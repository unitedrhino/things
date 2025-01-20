package msg

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPropertyLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogIndexLogic {
	return &PropertyLogIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PropertyLogIndexLogic) PropertyLogIndex(req *types.DeviceMsgPropertyLogIndexReq) (resp *types.DeviceMsgPropertyIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.PropertyLogIndex(l.ctx, utils.Copy[dm.PropertyLogIndexReq](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceData req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return utils.Copy[types.DeviceMsgPropertyIndexResp](dmResp), nil
}
