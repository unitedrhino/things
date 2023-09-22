package common

import (
	"context"
	"github.com/gorilla/websocket"
	ws "github.com/i-Things/things/shared/websocket"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebsocketLogic {
	return &WebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebsocketLogic) InitWebsocketConn(r *http.Request, conn *websocket.Conn) {

	//创建ws连接
	wsClient := ws.NewConn(l.svcCtx.Ws, r, conn)
	//开启读取进程
	go wsClient.StartRead()
	//开启发送进程
	go wsClient.StartWrite()
	return
}
