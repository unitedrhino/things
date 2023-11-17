package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"
	"github.com/zeromicro/go-zero/core/trace"
	"time"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
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

// 调用设备行为
func (l *SendActionLogic) SendAction(in *dm.SendActionReq) (*dm.SendActionResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	if err := CheckIsOnline(l.ctx, l.svcCtx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}); err != nil {
		return nil, err
	}

	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}

	param := map[string]any{}
	err = utils.Unmarshal([]byte(in.InputParams), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail("SendAction InputParams not right:", in.InputParams)
	}

	MsgToken := trace.TraceIDFromContext(l.ctx)

	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:    deviceMsg.Action,
			MsgToken:  MsgToken,
			Timestamp: time.Now().UnixMilli(),
		},
		ActionID: in.ActionID,
		Params:   param,
	}
	err = req.FmtReqParam(l.schema, schema.ParamActionInput)
	if err != nil {
		return nil, err
	}

	payload, _ := json.Marshal(req)
	reqMsg := deviceMsg.PublishMsg{
		Handle:     devices.Thing,
		Type:       msgThing.TypeAction,
		Payload:    payload,
		Timestamp:  time.Now().UnixMilli(),
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		Explain:    ToSendOptionDo(in.Option).String(),
	}
	err = cache.SetDeviceMsg(l.ctx, l.svcCtx.Cache, deviceMsg.ReqMsg, &reqMsg, req.MsgToken)
	if err != nil {
		return nil, err
	}

	if in.IsAsync { //如果是异步获取 处理结果暂不关注
		err := l.svcCtx.PubDev.PublishToDev(l.ctx, &reqMsg)
		if err != nil {
			return nil, err
		}
		if in.Option != nil {
			payload, _ := json.Marshal(reqMsg)
			_, err := l.svcCtx.TimedM.TaskSend(l.ctx, &timedjob.TaskSendReq{
				GroupCode: def.TimedIThingsQueueGroupCode,
				Code:      "disvr-action-check-delay",
				Option: &timedjob.TaskSendOption{
					ProcessIn: in.Option.RequestTimeout,
					Timeout:   in.Option.TimeoutToFail,
				},
				ParamQueue: &timedjob.TaskParamQueue{
					Topic:   topics.DiActionCheckDelay,
					Payload: string(payload),
				},
			})
			if err != nil {
				l.Errorf("TaskSend err:%v", err)
			}
		}
		return &dm.SendActionResp{
			MsgToken: req.MsgToken,
		}, nil
	}
	resp, err := l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, &reqMsg, func(payload []byte) bool {
		var dresp msgThing.Resp
		err = utils.Unmarshal(payload, &dresp)
		if err != nil { //如果是没法解析的说明不是需要的包,直接跳过即可
			return false
		}
		if dresp.MsgToken != req.MsgToken { //不是该请求的回复.跳过
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	var dresp msgThing.Resp
	err = utils.Unmarshal(resp, &dresp)
	if err != nil {
		return nil, err
	}
	respParam, err := json.Marshal(dresp.Data)
	if err != nil {
		return nil, errors.RespParam.AddDetailf("SendAction get device resp not right:%+v", dresp.Data)
	}
	return &dm.SendActionResp{
		MsgToken:     dresp.MsgToken,
		Msg:          dresp.Msg,
		Code:         dresp.Code,
		OutputParams: string(respParam),
	}, nil
}
