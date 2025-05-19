package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/timed/timedjobsvr/pb/timedjob"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/events/topics"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/cache"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ActionSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	schema *schema.Model
	logx.Logger
}

func NewActionSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionSendLogic {
	return &ActionSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ActionSendLogic) initMsg(dev devices.Core) error {
	var err error
	l.schema, err = l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, dev)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

// 调用设备行为
func (l *ActionSendLogic) ActionSend(in *dm.ActionSendReq) (ret *dm.ActionSendResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	_, err = logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthReadWrite, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, map[string]any{in.ActionID: struct{}{}})
	if err != nil {
		return nil, err
	}
	var protocolCode string
	if protocolCode, err = CheckIsOnline(l.ctx, l.svcCtx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}); err != nil {
		return nil, err
	}

	err = l.initMsg(devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	if err != nil {
		return nil, err
	}

	param := map[string]any{}
	if len(in.InputParams) > 0 {
		err = utils.Unmarshal([]byte(in.InputParams), &param)
		if err != nil {
			return nil, errors.Parameter.AddDetail("ActionSend InputParams not right:", in.InputParams)
		}
	}
	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:   deviceMsg.Action,
			MsgToken: devices.GenMsgToken(l.ctx, l.svcCtx.NodeID),
			//Timestamp: time.Now().UnixMilli(),
		},
		ActionID: in.ActionID,
		Params:   param,
	}
	params, err := req.VerifyReqParam(l.schema, schema.ParamActionInput)
	if err != nil {
		return nil, err
	}
	req.Params, err = msgThing.ToVal(params)
	if err != nil {
		return nil, err
	}
	defer func() {
		ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
			uc := ctxs.GetUserCtx(l.ctx)
			var content = map[string]any{}
			content["req"] = params
			content["userID"] = uc.UserID
			contentStr, _ := json.Marshal(params)
			di, err := l.svcCtx.DeviceCache.GetData(ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
			if err != nil {
				l.Error(err)
				return
			}
			_ = l.svcCtx.SendRepo.Insert(ctx, &deviceLog.Send{
				TenantCode:  dataType.TenantCode(di.TenantCode),
				ProjectID:   dataType.ProjectID(di.ProjectID),
				AreaID:      dataType.AreaID(di.AreaID),
				AreaIDPath:  dataType.AreaIDPath(di.AreaIDPath),
				BelongGroup: utils.CopyMap3[def.IDsInfo](di.BelongGroup),
				ProductID:   in.ProductID,
				Account:     uc.Account,
				Action:      "actionSend",
				Timestamp:   time.Now(), // 操作时间
				DeviceName:  in.DeviceName,
				TraceID:     utils.TraceIdFromContext(ctx),
				UserID:      uc.UserID,
				DataID:      in.ActionID,
				Content:     string(contentStr),
				ResultCode:  errors.Fmt(err).GetCode(),
			})
		})
	}()
	payload, _ := json.Marshal(req)
	reqMsg := deviceMsg.PublishMsg{
		Handle:       devices.Thing,
		Type:         msgThing.TypeAction,
		Payload:      payload,
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    in.ProductID,
		DeviceName:   in.DeviceName,
		Explain:      ToSendOptionDo(in.Option).String(),
		ProtocolCode: protocolCode,
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
				GroupCode: def.TimedUnitedRhinoQueueGroupCode,
				Code:      "disvr-action-check-delay",
				Option: &timedjob.TaskSendOption{
					ProcessIn: in.Option.RequestTimeout,
					Timeout:   in.Option.TimeoutToFail,
					TaskID:    req.MsgToken,
				},
				ParamQueue: &timedjob.TaskParamQueue{
					Topic:   topics.DmActionCheckDelay,
					Payload: string(payload),
				},
			})
			if err != nil {
				l.Errorf("TaskSend err:%v", err)
			}
		}
		return &dm.ActionSendResp{
			MsgToken: req.MsgToken,
		}, nil
	}
	resp, err := l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, &reqMsg, 0, func(payload []byte) bool {
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
		return nil, errors.RespParam.AddDetailf("ActionSend get device resp not right:%+v", dresp.Data)
	}
	return &dm.ActionSendResp{
		MsgToken:     dresp.MsgToken,
		Msg:          dresp.Msg,
		Code:         dresp.Code,
		OutputParams: string(respParam),
	}, nil
}
