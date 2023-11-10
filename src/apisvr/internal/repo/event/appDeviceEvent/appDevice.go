package appDeviceEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	ws "github.com/i-Things/things/shared/websocket"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
	"net/http"
)

type AppDeviceHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewAppDeviceHandle(ctx context.Context, svcCtx *svc.ServiceContext) *AppDeviceHandle {
	return &AppDeviceHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (a *AppDeviceHandle) DeviceEventReport(in *application.EventReport) error {
	a.Infof("%s req=%v", utils.FuncName(), in)
	return nil
}

func (a *AppDeviceHandle) DevicePropertyReport(in *application.PropertyReport) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceReportThingPropertyDevice, in.Device.ProductID, in.Device.DeviceName)
	MsgToken := trace.TraceIDFromContext(a.ctx)
	param := map[string]interface{}{
		in.Identifier: in.Param,
	}
	data, _ := json.Marshal(param)
	body := ws.WsBody{
		Type: ws.Pub,
		Path: topic,
		Body: string(data),
	}
	body.Handler = make(http.Header)
	body.Handler.Set("Traceparent", MsgToken)
	ws.SendSub(ws.WsResp{
		StatusCode: http.StatusOK,
		WsBody:     body,
	})
	return nil
}

func (a *AppDeviceHandle) DeviceStatusConnected(in *application.ConnectMsg) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusConnected, in.Device.ProductID, in.Device.DeviceName)
	MsgToken := trace.TraceIDFromContext(a.ctx)
	body := ws.WsBody{
		Type: ws.Pub,
		Path: topic,
		Body: "connected",
	}
	body.Handler = make(http.Header)
	body.Handler.Set("Traceparent", MsgToken)
	ws.SendSub(ws.WsResp{
		StatusCode: http.StatusOK,
		WsBody:     body,
	})
	return nil
}

func (a *AppDeviceHandle) DeviceStatusDisConnected(in *application.ConnectMsg) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusDisConnected, in.Device.ProductID, in.Device.DeviceName)
	MsgToken := trace.TraceIDFromContext(a.ctx)
	body := ws.WsBody{
		Type: ws.Pub,
		Path: topic,
		Body: "disconnected",
	}
	body.Handler = make(http.Header)
	body.Handler.Set("Traceparent", MsgToken)
	ws.SendSub(ws.WsResp{
		StatusCode: http.StatusOK,
		WsBody:     body,
	})
	return nil
}
