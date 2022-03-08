package dm

import (
	"net/http"

	"github.com/go-things/things/src/webapi/internal/logic/backgrand/dm"
	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetDevicePropertyStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetDevicePropertyStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dm.NewGetDevicePropertyStatusLogic(r.Context(), svcCtx)
		resp, err := l.GetDevicePropertyStatus(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
