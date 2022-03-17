package dc

import (
	"github.com/i-Things/things/src/webapi/internal/logic/front/dc"
	"net/http"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetGroupInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGroupInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dc.NewGetGroupInfoLogic(r.Context(), ctx)
		resp, err := l.GetGroupInfo(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
