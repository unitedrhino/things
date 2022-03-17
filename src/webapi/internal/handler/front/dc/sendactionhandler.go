package dc

import (
	"github.com/i-Things/things/src/webapi/internal/logic/front/dc"
	"net/http"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendActionHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendDcActionReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dc.NewSendActionLogic(r.Context(), ctx)
		resp, err := l.SendAction(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
