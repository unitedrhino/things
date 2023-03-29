package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"time"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendActionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	schema *schema.Model
	logx.Logger
}

func NewSendActionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendActionLogic {
	return &SendActionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *SendActionLogic) initMsg(productID string) error {
	var err error
	l.schema, err = l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

func (l *SendActionLogic) SendAction(in *di.SendActionReq) (*di.SendActionResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}
	param := map[string]any{}
	err = utils.Unmarshal([]byte(in.InputParams), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail("SendAction InputParams not right:", in.InputParams)
	}
	uuid, err := uuid.GenerateUUID()
	if err != nil {
		l.Errorf("%s.GenerateUUID err:%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:      deviceMsg.Action,
			ClientToken: uuid,
			Timestamp:   time.Now().UnixMilli(),
		},
		ActionID: in.ActionID,
		Params:   param}
	_, err = req.VerifyReqParam(l.schema, schema.ParamActionInput)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.MsgThingRepo.SetReq(l.ctx, msgThing.TypeAction,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, &req)
	if err != nil {
		return nil, err
	}
	pubTopic := fmt.Sprintf("$thing/down/action/%s/%s", in.ProductID, in.DeviceName)
	subTopic := fmt.Sprintf("$thing/up/action/%s/%s", in.ProductID, in.DeviceName)

	if in.IsAsync { //如果是异步获取 处理结果暂不关注
		payload, _ := json.Marshal(req)
		err := l.svcCtx.PubDev.PublishToDev(l.ctx, pubTopic, payload)
		if err != nil {
			return nil, err
		}
		return &di.SendActionResp{
			ClientToken: req.ClientToken,
		}, nil
	}
	resp, err := l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, pubTopic, subTopic, &req, in.ProductID, in.DeviceName)
	if err != nil {
		return nil, err
	}
	respParam, err := json.Marshal(resp.Response)
	if err != nil {
		return nil, errors.RespParam.AddDetailf("SendAction get device resp not right:%+v", resp.Response)
	}
	return &di.SendActionResp{
		ClientToken:  resp.ClientToken,
		Status:       resp.Status,
		Code:         resp.Code,
		OutputParams: string(respParam),
	}, nil
}
