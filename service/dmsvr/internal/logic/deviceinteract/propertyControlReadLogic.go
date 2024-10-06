package deviceinteractlogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/cache"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

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
func (l *PropertyControlReadLogic) PropertyControlRead(in *dm.RespReadReq) (*dm.PropertyControlSendResp, error) {
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cache.GetDeviceMsg[msgThing.Resp](l.ctx, l.svcCtx.Cache, deviceMsg.RespMsg, devices.Thing, msgThing.TypeProperty,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, in.MsgToken)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	if resp == nil {
		return nil, errors.NotFind
	}
	return &dm.PropertyControlSendResp{
		MsgToken: resp.MsgToken,
		Msg:      resp.Msg,
		Code:     resp.Code,
	}, nil
}
