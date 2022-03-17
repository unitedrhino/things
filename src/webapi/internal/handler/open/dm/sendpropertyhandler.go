package dm

import (
	"net/http"

	"github.com/i-Things/things/src/webapi/internal/logic/open/dm"
	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendPropertyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendDmPropertyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dm.NewSendPropertyLogic(r.Context(), svcCtx)
		resp, err := l.SendProperty(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
