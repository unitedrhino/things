package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/core/service/syssvr/sysExport"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/application"
	"gitee.com/unitedrhino/share/domain/deviceMsg"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/shadow"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	devicemanagelogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/devicemanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/cache"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	l.model, err = l.svcCtx.SchemaRepo.GetData(l.ctx, productID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	return nil
}

// 调用设备属性
func (l *PropertyControlSendLogic) PropertyControlSend(in *dm.PropertyControlSendReq) (ret *dm.PropertyControlSendResp, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), in)
	var isOnline = true
	var protocolCode string
	var dev = devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	if protocolCode, err = CheckIsOnline(l.ctx, l.svcCtx, dev); err != nil { //如果是不启用设备影子的模式则直接返回
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
	param, err = logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthReadWrite, dev, param)
	if err != nil {
		return nil, err
	}
	MsgToken := devices.GenMsgToken(l.ctx, l.svcCtx.NodeID)

	req := msgThing.Req{
		CommonMsg: deviceMsg.CommonMsg{
			Method:   deviceMsg.Control,
			MsgToken: MsgToken,
			//Timestamp: time.Now().UnixMilli(),
		},
		Params: param,
	}
	params, err := req.VerifyReqParam(l.model, schema.ParamProperty)
	if err != nil {
		return nil, err
	}
	if len(params) == 0 {
		l.Infof("控制的属性在设备中都不存在,req:%v", utils.Fmt(in))
		return &dm.PropertyControlSendResp{Code: errors.OK.Code, Msg: errors.OK.AddMsg("该设备无控制的属性,忽略").GetMsg()}, nil
	}
	req.Params, err = msgThing.ToVal(params)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil && ret.Code == errors.OK.GetCode() && in.WithProfile != nil && in.WithProfile.Code != "" {
			_, err = devicemanagelogic.NewDeviceProfileUpdateLogic(l.ctx, l.svcCtx).DeviceProfileUpdate(&dm.DeviceProfile{
				Device: &dm.DeviceCore{
					ProductID:  in.ProductID,
					DeviceName: in.DeviceName,
				},
				Code:   in.WithProfile.Code,
				Params: in.WithProfile.Params,
			})
		}
	}()
	if utils.SliceIn(in.ShadowControl, shadow.ControlOnlyCloud, shadow.ControlOnlyCloudWithLog) {
		for k, v := range req.Params {
			appMsg := application.PropertyReport{
				Device: dev, Timestamp: time.Now().UnixMilli(),
				Identifier: k, Param: v,
			}
			//应用事件通知-设备物模型属性上报通知 ↓↓↓
			err := l.svcCtx.PubApp.DeviceThingPropertyReport(l.ctx, appMsg)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("%s.DeviceThingPropertyReport  identifier:%v, param:%v,err:%v", utils.FuncName(), k, param, err)
			}
			err = l.svcCtx.WebHook.Publish(l.svcCtx.WithDeviceTenant(l.ctx, dev), sysExport.CodeDmDevicePropertyReport, appMsg)
			if err != nil {
				l.Error(err)
			}
			err = l.svcCtx.UserSubscribe.Publish(l.ctx, def.UserSubscribeDevicePropertyReport, appMsg, map[string]any{
				"productID":  in.ProductID,
				"deviceName": in.DeviceName,
				"identifier": k,
			}, map[string]any{
				"productID":  dev.ProductID,
				"deviceName": dev.DeviceName,
			})
		}
	}
	if in.ShadowControl == shadow.ControlOnlyCloud {
		//插入多条设备物模型属性数据
		err = l.svcCtx.SchemaManaRepo.InsertPropertiesData(l.ctx, l.model, in.ProductID, in.DeviceName, params, time.Now())
		if err != nil {
			l.Errorf("%s.InsertPropertyData err=%+v", utils.FuncName(), err)
			return nil, err
		}
		return &dm.PropertyControlSendResp{Code: errors.OK.Code, Msg: errors.OK.AddMsg("只修改云端值").GetMsg()}, nil
	}
	defer func() {
		ctxs.GoNewCtx(l.ctx, func(ctx context.Context) {
			uc := ctxs.GetUserCtxNoNil(l.ctx)
			account := uc.Account
			if account == "" && uc.UserID <= def.RootNode {
				account = "系统控制"
			}
			for dataID, content := range param {
				_ = l.svcCtx.SendRepo.Insert(ctx, &deviceLog.Send{
					ProductID:  in.ProductID,
					Action:     "propertyControlSend",
					Timestamp:  time.Now(), // 操作时间
					DeviceName: in.DeviceName,
					TraceID:    utils.TraceIdFromContext(ctx),
					UserID:     uc.UserID,
					DataID:     dataID,
					Account:    account,
					Content:    utils.Fmt(content),
					ResultCode: errors.Fmt(err).GetCode(),
				})
			}
		})
	}()
	if in.ShadowControl == shadow.ControlOnlyCloudWithLog {
		//插入多条设备物模型属性数据
		err = l.svcCtx.SchemaManaRepo.InsertPropertiesData(l.ctx, l.model, in.ProductID, in.DeviceName, params, time.Now())
		if err != nil {
			l.Errorf("%s.InsertPropertyData err=%+v", utils.FuncName(), err)
			return nil, err
		}
		return &dm.PropertyControlSendResp{Code: errors.OK.Code, Msg: errors.OK.AddMsg("只修改云端值及记录操作").GetMsg()}, nil
	}
	if in.ShadowControl == shadow.ControlOnly || (!isOnline && in.ShadowControl == shadow.ControlAuto) {
		//设备影子模式
		err = shadow.CheckEnableShadow(param, l.model)
		if err != nil {
			if !isOnline && in.ShadowControl == shadow.ControlAuto { //如果是自动且不在线的模式
				err = errors.NotOnline
			}
			return nil, err
		}
		err = relationDB.NewShadowRepo(l.ctx).MultiUpdate(l.ctx, shadow.NewInfo(in.ProductID, in.DeviceName, param))
		if err != nil {
			return nil, err
		}
		return &dm.PropertyControlSendResp{Code: errors.OK.Code, Msg: errors.OK.AddMsg("影子模式").GetMsg()}, nil
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
	resp, err = l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, &reqMsg, time.Duration(in.SyncTimeout)*time.Second, func(payload []byte) bool {
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
		return nil, errors.Fmt(err).WithMsg("指令发送失败")
	}

	var dresp msgThing.Resp
	err = utils.Unmarshal(resp, &dresp)
	if err != nil {
		return nil, err
	}
	if dresp.Code != errors.OK.GetCode() {
		if dresp.Msg != "" {
			err = errors.DeviceResp.AddMsg(dresp.Msg)
		} else {
			err = errors.DeviceResp
		}
		err = errors.Fmt(err).WithMsg("指令发送失败")
	}
	ret = &dm.PropertyControlSendResp{
		MsgToken: dresp.MsgToken,
		Msg:      dresp.Msg,
		Code:     dresp.Code,
	}
	return ret, err
}
