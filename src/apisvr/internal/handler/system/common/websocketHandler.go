package common

import (
	"github.com/gorilla/websocket"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/common"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

func WebsocketHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			// 读取存储空间大小
			ReadBufferSize: 1024,
			// 写入存储空间大小
			WriteBufferSize: 1024,
			// 允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		//ws连接失败
		if err != nil {
			result.Http(w, r, nil, err)
			logx.WithContext(r.Context()).Error("[ws]连接失败", "RemoteAddr:", r.RemoteAddr, "err", err)
			return
		}
		l := common.NewWebsocketLogic(r.Context(), svcCtx)
		l.InitWebsocketConn(r, conn)
	}
}
