package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgGateway"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/logic"
	"time"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayNotifyBindSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayNotifyBindSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayNotifyBindSendLogic {
	return &GatewayNotifyBindSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知网关绑定子设备
func (l *GatewayNotifyBindSendLogic) GatewayNotifyBindSend(in *dm.GatewayNotifyBindSendReq) (*dm.Empty, error) {
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}, nil)
	if err != nil {
		return nil, err
	}
	var protocolCode string
	if protocolCode, err = CheckIsOnline(l.ctx, l.svcCtx, devices.Core{
		ProductID:  in.Gateway.ProductID,
		DeviceName: in.Gateway.DeviceName,
	}); err != nil {
		return nil, err
	}
	MsgToken := devices.GenMsgToken(l.ctx, l.svcCtx.NodeID)
	req := msgGateway.Msg{
		CommonMsg: deviceMsg.CommonMsg{
			Method:   deviceMsg.NotifyBind,
			MsgToken: MsgToken,
			//Timestamp: time.Now().UnixMilli(),
		},
		Payload: &msgGateway.GatewayPayload{
			Devices: nil,
		},
	}
	for _, v := range in.SubDevices {
		req.Payload.Devices = append(req.Payload.Devices, &msgGateway.Device{
			ProductID:  v.ProductID,
			DeviceName: v.DeviceName,
		})
	}
	payload, _ := json.Marshal(req)
	reqMsg := deviceMsg.PublishMsg{
		Handle:       devices.Gateway,
		Type:         msgGateway.TypeTopo,
		Payload:      payload,
		Timestamp:    time.Now().UnixMilli(),
		ProductID:    in.Gateway.ProductID,
		DeviceName:   in.Gateway.DeviceName,
		ProtocolCode: protocolCode,
	}
	var resp []byte
	resp, err = l.svcCtx.PubDev.ReqToDeviceSync(l.ctx, &reqMsg, time.Second*10, func(payload []byte) bool {
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
	var dresp msgGateway.Msg
	err = utils.Unmarshal(resp, &dresp)
	if err != nil {
		return nil, err
	}
	if dresp.Code != errors.OK.GetCode() {
		return nil, errors.DeviceResp.AddMsgf("设备返回错误: msg:%v code:%v", dresp.Msg, dresp.Code)
	}
	return &dm.Empty{}, err
}
