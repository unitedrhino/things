package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/disvr/internal/domain/shadow"
	"github.com/i-Things/things/src/disvr/internal/repo/cache"
	"github.com/i-Things/things/src/disvr/internal/repo/relationDB"
	"time"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendPropertyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	model *schema.Model
}

func NewSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendPropertyLogic {
	return &SendPropertyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendPropertyLogic) initMsg(productID string) error {
	var err error
	l.model, err = l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

func (l *SendPropertyLogic) SendProperty(in *di.SendPropertyReq) (*di.SendPropertyResp, error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	var isOnline = true
	if err := checkIsOnline(l.ctx, l.svcCtx, devices.Core{
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

	err := l.initMsg(in.ProductID)
	if err != nil {
		return nil, err
	}

	param := map[string]any{}
	err = utils.Unmarshal([]byte(in.Data), &param)
	if err != nil {
		return nil, errors.Parameter.AddDetail(
			"SendProperty data not right:", in.Data)
	}
	clientToken, err := uuid.GenerateUUID()
	if err != nil {
		l.Errorf("%s.GenerateUUID err:%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}

	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:      deviceMsg.Control,
			ClientToken: clientToken,
			Timestamp:   time.Now().UnixMilli(),
		},
		Params: param,
	}
	_, err = req.VerifyReqParam(l.model, schema.ParamProperty)
	if err != nil {
		return nil, err
	}
	if in.ShadowControl == shadow.ControlOnly || (!isOnline && in.ShadowControl == shadow.ControlAuto) {
		//设备影子模式
		err := shadow.CheckEnableShadow(param, l.model)
		if err != nil {
			return nil, err
		}
		err = relationDB.NewShadowRepo(l.ctx).MultiUpdate(l.ctx, shadow.NewInfo(in.ProductID, in.DeviceName, param))
		if err != nil {
			return nil, err
		}
		return &di.SendPropertyResp{}, nil
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
	err = cache.SetDeviceMsg(l.ctx, l.svcCtx.Cache, deviceMsg.ReqMsg, &reqMsg, req.ClientToken)
	if err != nil {
		return nil, err
	}

	if in.IsAsync { //如果是异步获取 处理结果暂不关注
		err := l.svcCtx.PubDev.PublishToDev(l.ctx, &reqMsg)
		if err != nil {
			return nil, err
		}
		return &di.SendPropertyResp{
			ClientToken: req.ClientToken,
		}, nil
	}

	resp, err := l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, &reqMsg, func(payload []byte) bool {
		var dresp msgThing.Resp
		err = utils.Unmarshal(payload, &dresp)
		if err != nil { //如果是没法解析的说明不是需要的包,直接跳过即可
			return false
		}
		if dresp.ClientToken != req.ClientToken { //不是该请求的回复.跳过
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

	return &di.SendPropertyResp{
		ClientToken: dresp.ClientToken,
		Status:      dresp.Status,
		Code:        dresp.Code,
	}, nil
}
