package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"github.com/i-Things/things/service/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"time"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ActionRespLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	schema *schema.Model
}

func NewActionRespLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionRespLogic {
	return &ActionRespLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ActionRespLogic) initMsg(productID string) error {
	var err error
	l.schema, err = l.svcCtx.SchemaRepo.GetData(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

// 回复调用设备行为
func (l *ActionRespLogic) ActionResp(in *dm.ActionRespReq) (*dm.Empty, error) {
	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}
	req, err := cache.GetDeviceMsg[msgThing.Req](l.ctx, l.svcCtx.Cache, deviceMsg.ReqMsg, devices.Thing, deviceMsg.Action,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName},
		in.MsgToken)
	if req == nil || err != nil {
		return nil, err
	}
	_, err = logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthReadWrite, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, map[string]any{req.ActionID: struct{}{}})
	if err != nil {
		return nil, err
	}
	resp := msgThing.Resp{
		CommonMsg: deviceMsg.CommonMsg{
			Method:   deviceMsg.ActionReply,
			MsgToken: in.MsgToken,
			//Timestamp: time.Now().UnixMilli(),
			Msg:  in.Msg,
			Code: in.Code,
		},
		ActionID: req.ActionID,
	}
	if resp.Code == 0 {
		resp.Code = errors.OK.Code
		resp.Msg = errors.OK.GetMsg()
	}
	if resp.Code == errors.OK.GetCode() {
		param := map[string]any{}
		err = utils.Unmarshal([]byte(in.OutputParams), &param)
		if err != nil {
			return nil, errors.Parameter.AddDetail("SendAction InputParams not right:", in.OutputParams)
		}
		resp.Data = param
		err = resp.FmtRespParam(l.schema, req.ActionID, schema.ParamActionOutput)
		if err != nil {
			return nil, err
		}
	}

	payload, _ := json.Marshal(resp)
	reqMsg := deviceMsg.PublishMsg{
		Handle:     devices.Thing,
		Type:       msgThing.TypeAction,
		Payload:    payload,
		Timestamp:  time.Now().UnixMilli(),
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	err = l.svcCtx.PubDev.PublishToDev(l.ctx, &reqMsg)
	if err != nil {
		return nil, err
	}

	return &dm.Empty{}, nil
}
