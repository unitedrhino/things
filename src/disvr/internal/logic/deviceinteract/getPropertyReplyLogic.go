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
	"github.com/zeromicro/go-zero/core/trace"
	"time"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPropertyReplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	schema *schema.Model
	logx.Logger
}

func NewGetPropertyReplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPropertyReplyLogic {
	return &GetPropertyReplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *GetPropertyReplyLogic) initMsg(productID string) error {
	var err error
	l.schema, err = l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

// 请求设备获取设备最新属性
func (l *GetPropertyReplyLogic) GetPropertyReply(in *di.GetPropertyReplyReq) (*di.GetPropertyReplyResp, error) {
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

	clientToken := trace.TraceIDFromContext(l.ctx)

	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:      deviceMsg.GetReport,
			ClientToken: clientToken,
			Timestamp:   time.Now().UnixMilli(),
		},
		Identifiers: in.DataIDs,
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

	var dresp msgThing.Req
	err = utils.Unmarshal(resp, &dresp)
	if err != nil {
		return nil, err
	}
	var params []byte
	if len(dresp.Params) > 0 {
		params, _ = json.Marshal(dresp.Params)
	}
	return &di.GetPropertyReplyResp{
		ClientToken: dresp.ClientToken,
		Status:      dresp.Status,
		Code:        dresp.Code,
		Timestamp:   dresp.GetTimeStamp(time.Now().UnixMilli()).UnixMilli(),
		Params:      string(params),
	}, nil
}
