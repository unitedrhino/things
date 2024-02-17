package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/service/dmsvr/internal/domain/shadow"
	"github.com/i-Things/things/service/dmsvr/internal/repo/cache"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyControlSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	model *schema.Model
}

func NewPropertyControlSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyControlSendLogic {
	return &PropertyControlSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *PropertyControlSendLogic) initMsg(productID string) error {
	var err error
	l.model, err = l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

// 调用设备属性
func (l *PropertyControlSendLogic) PropertyControlSend(in *dm.PropertyControlSendReq) (ret *dm.PropertyControlSendResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	var isOnline = true
	if err := CheckIsOnline(l.ctx, l.svcCtx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}); err != nil { //如果是不启用设备影子的模式则直接返回
		if errors.Is(err, errors.NotOnline) {
			isOnline = false
			if in.ShadowControl == shadow.ControlNo {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	err = l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}

	param := map[string]any{}
	err = utils.Unmarshal([]byte(in.Data), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail(
			"SendProperty data not right:", in.Data)
	}
	MsgToken := devices.GenMsgToken(l.ctx)

	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:    deviceMsg.Control,
			MsgToken:  MsgToken,
			Timestamp: time.Now().UnixMilli(),
		},
		Params: param,
	}
	params, err := req.VerifyReqParam(l.model, schema.ParamProperty)
	if err != nil {
		return nil, err
	}
	err = req.FmtReqParam(l.model, schema.ParamProperty)
	if err != nil {
		return nil, err
	}
	req.Params, err = msgThing.ToVal(params)
	if err != nil {
		return nil, err
	}
	defer func() {
		utils.GoNewCtx(l.ctx, func(ctx context.Context) {
			uc := ctxs.GetUserCtx(l.ctx)
			var content = map[string]any{}
			content["req"] = params
			content["userID"] = uc.UserID
			contentStr, _ := json.Marshal(content)
			_ = l.svcCtx.HubLogRepo.Insert(ctx, &msgHubLog.HubLog{
				ProductID:  in.ProductID,
				Action:     "propertyControlSend",
				Timestamp:  time.Now(), // 操作时间
				DeviceName: in.DeviceName,
				TranceID:   utils.TraceIdFromContext(ctx),
				RequestID:  MsgToken,
				Content:    string(contentStr),
				ResultType: errors.Fmt(err).GetCode(),
			})
		})
	}()
	if in.ShadowControl == shadow.ControlOnly || (!isOnline && in.ShadowControl == shadow.ControlAuto) {
		//设备影子模式
		err = shadow.CheckEnableShadow(param, l.model)
		if err != nil {
			return nil, err
		}
		err = relationDB.NewShadowRepo(l.ctx).MultiUpdate(l.ctx, shadow.NewInfo(in.ProductID, in.DeviceName, param))
		if err != nil {
			return nil, err
		}
		return &dm.PropertyControlSendResp{}, nil
	}

	payload, _ := json.Marshal(req)
	reqMsg := deviceMsg.PublishMsg{
		Handle:     devices.Thing,
		Type:       msgThing.TypeProperty,
		Payload:    payload,
		Timestamp:  time.Now().UnixMilli(),
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
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
		return &dm.PropertyControlSendResp{
			MsgToken: req.MsgToken,
		}, nil
	}
	var resp []byte
	resp, err = l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, &reqMsg, func(payload []byte) bool {
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

	return &dm.PropertyControlSendResp{
		MsgToken: dresp.MsgToken,
		Msg:      dresp.Msg,
		Code:     dresp.Code,
	}, nil
}
