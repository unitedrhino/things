package dm

import (
	"github.com/go-things/things/src/webapi/internal/logic/backgrand/dm"
	"net/http"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetDeviceInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDeviceInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dm.NewGetDeviceInfoLogic(r.Context(), ctx)
		resp, err := l.GetDeviceInfo(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
