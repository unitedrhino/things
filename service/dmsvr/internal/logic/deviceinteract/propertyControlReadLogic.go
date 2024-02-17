package deviceinteractlogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/service/dmsvr/internal/repo/cache"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyControlReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyControlReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyControlReadLogic {
	return &PropertyControlReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取异步调用设备属性的结果
func (l *PropertyControlReadLogic) PropertyControlRead(in *dm.RespReadReq) (*dm.SendPropertyControlResp, error) {
	resp, err := cache.GetDeviceMsg[msgThing.Resp](l.ctx, l.svcCtx.Cache, deviceMsg.RespMsg, devices.Thing, msgThing.TypeProperty,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, in.MsgToken)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	if resp == nil {
		return nil, errors.NotFind
	}
	return &dm.SendPropertyControlResp{
		MsgToken: resp.MsgToken,
		Msg:      resp.Msg,
		Code:     resp.Code,
	}, nil
}
