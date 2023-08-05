package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/disvr/internal/repo/cache"
	"time"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type RespActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	schema *schema.Model
}

func NewRespActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RespActionLogic {
	return &RespActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *RespActionLogic) initMsg(productID string) error {
	var err error
	l.schema, err = l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

// 回复调用设备行为
func (l *RespActionLogic) RespAction(in *di.RespActionReq) (*di.Response, error) {
	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}
	req, err := cache.GetDeviceMsg[msgThing.Req](l.ctx, l.svcCtx.Cache, deviceMsg.ReqMsg, devices.Thing, deviceMsg.Action,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName},
		in.ClientToken)
	if req == nil || err != nil {
		return nil, err
	}
	param := map[string]any{}
	err = utils.Unmarshal([]byte(in.OutputParams), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail("SendAction InputParams not right:", in.OutputParams)
	}
	resp := msgThing.Resp{
		CommonMsg: deviceMsg.CommonMsg{
			Method:      deviceMsg.Action,
			ClientToken: in.ClientToken,
			Timestamp:   time.Now().UnixMilli(),
			Status:      in.Status,
			Code:        in.Code,
		},
		Response: param,
	}
	if resp.Code == 0 {
		resp.Code = errors.OK.Code
		resp.Status = errors.OK.Msg
	}

	err = resp.FmtRespParam(l.schema, req.ActionID, schema.ParamActionOutput)
	if err != nil {
		return nil, err
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

	return &di.Response{}, nil
}
