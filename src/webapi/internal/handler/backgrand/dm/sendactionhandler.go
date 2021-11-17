package dm

import (
	"gitee.com/godLei6/things/src/webapi/internal/logic/backgrand/dm"
	"net/http"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func SendActionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendDmActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dm.NewSendActionLogic(r.Context(), ctx)
		resp, err := l.SendAction(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
