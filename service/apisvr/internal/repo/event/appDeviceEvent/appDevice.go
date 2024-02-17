package appDeviceEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/domain/application"
	"gitee.com/i-Things/share/events/topics"
	"gitee.com/i-Things/share/utils"
	ws "gitee.com/i-Things/share/websocket"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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
	param := map[string]any{
		in.Identifier: in.Param,
	}
	data, _ := json.Marshal(param)
	body := ws.WsBody{
		Type: ws.Pub,
		Path: topic,
		Body: string(data),
	}
	ws.SendSub(a.ctx, ws.WsResp{
		StatusCode: http.StatusOK,
		WsBody:     body,
	})
	return nil
}

func (a *AppDeviceHandle) DeviceStatusConnected(in *application.ConnectMsg) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusConnected, in.Device.ProductID, in.Device.DeviceName)
	body := ws.WsBody{
		Type: ws.Pub,
		Path: topic,
		Body: def.ConnectedStatus,
	}
	ws.SendSub(a.ctx, ws.WsResp{
		StatusCode: http.StatusOK,
		WsBody:     body,
	})
	return nil
}

func (a *AppDeviceHandle) DeviceStatusDisConnected(in *application.ConnectMsg) error {
	topic := fmt.Sprintf(topics.ApplicationDeviceStatusDisConnected, in.Device.ProductID, in.Device.DeviceName)
	body := ws.WsBody{
		Type: ws.Pub,
		Path: topic,
		Body: def.DisConnectedStatus,
	}
	ws.SendSub(a.ctx,
		ws.WsResp{StatusCode: http.StatusOK,
			WsBody: body,
		})
	return nil
}
