package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/cache"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"strings"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyGetReportSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	schema *schema.Model
}

func NewPropertyGetReportSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyGetReportSendLogic {
	return &PropertyGetReportSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *PropertyGetReportSendLogic) initMsg(dev devices.Core) error {
	var err error
	l.schema, err = l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, dev)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

// 请求设备获取设备最新属性
func (l *PropertyGetReportSendLogic) PropertyGetReportSend(in *dm.PropertyGetReportSendReq) (ret *dm.PropertyGetReportSendResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	_, err = logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
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

	MsgToken := devices.GenMsgToken(l.ctx, l.svcCtx.NodeID)

	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:   deviceMsg.GetReport,
			MsgToken: MsgToken,
			//Timestamp: time.Now().UnixMilli(),
		},
		Identifiers: in.DataIDs,
	}

	payload, _ := json.Marshal(req)
	reqMsg := deviceMsg.PublishMsg{
		Handle:       devices.Thing,
		Type:         msgThing.TypeProperty,
		Payload:      payload,
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    in.ProductID,
		DeviceName:   in.DeviceName,
		ProtocolCode: protocolCode,
	}
	err = cache.SetDeviceMsg(l.ctx, l.svcCtx.Cache, deviceMsg.ReqMsg, &reqMsg, req.MsgToken)
	if err != nil {
		return nil, err
	}
	var resp []byte
	var params []byte
	defer func() {
		uc := ctxs.GetUserCtxNoNil(l.ctx)
		account := uc.Account
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
			ProductID:  in.ProductID,
			DeviceName: in.DeviceName,
		})
		if err != nil {
			l.Error(err)
			return
		}
		do := deviceLog.Send{
			TenantCode:  di.TenantCode,
			ProjectID:   di.ProjectID,
			AreaID:      di.AreaID,
			AreaIDPath:  di.AreaIDPath,
			BelongGroup: utils.CopyMap3[def.IDsInfo](di.BelongGroup),
			ProductID:   in.ProductID,
			Action:      "propertyGetReportSend",
			Timestamp:   time.Now(), // 操作时间
			DataID:      strings.Join(in.DataIDs, ","),
			DeviceName:  in.DeviceName,
			TraceID:     utils.TraceIdFromContext(l.ctx),
			UserID:      uc.UserID,
			Account:     account,
			Content:     string(params),
			ResultCode:  errors.Fmt(err).GetCode(),
		}

		_ = l.svcCtx.SendRepo.Insert(l.ctx, &do)
	}()
	resp, err = l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, &reqMsg, 0, func(payload []byte) bool {
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

	var dresp msgThing.Req
	err = utils.Unmarshal(resp, &dresp)
	if err != nil {
		return nil, err
	}
	if dresp.Data != nil {
		params, _ = json.Marshal(dresp.Data)
	}

	return &dm.PropertyGetReportSendResp{
		MsgToken:  dresp.MsgToken,
		Msg:       dresp.Msg,
		Code:      dresp.Code,
		Timestamp: dresp.GetTimeStamp(time.Now().UnixMilli()).UnixMilli(),
		Params:    string(params),
	}, nil
}
