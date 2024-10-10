package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceMsg"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/cache"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ActionReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionReadLogic {
	return &ActionReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取异步调用设备行为的结果
func (l *ActionReadLogic) ActionRead(in *dm.RespReadReq) (*dm.ActionSendResp, error) {
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
	if err != nil {
		return nil, err
	}
	resp, err := cache.GetDeviceMsg[msgThing.Resp](l.ctx, l.svcCtx.Cache, deviceMsg.RespMsg, devices.Thing, msgThing.TypeAction,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, in.MsgToken)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	if resp == nil {
		return nil, errors.NotFind
	}
	respParam, err := json.Marshal(resp.Data)
	if err != nil {
		return nil, errors.RespParam.AddDetailf("SendAction get device resp not right:%+v", resp.Data)
	}
	return &dm.ActionSendResp{
		MsgToken:     resp.MsgToken,
		Msg:          resp.Msg,
		Code:         resp.Code,
		OutputParams: string(respParam),
	}, nil
}
