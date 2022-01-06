package dc

import (
	"net/http"

	"github.com/go-things/things/src/webapi/internal/logic/backgrand/dc"
	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func BgGetGroupInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGroupInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dc.NewBgGetGroupInfoLogic(r.Context(), ctx)
		resp, err := l.BgGetGroupInfo(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
