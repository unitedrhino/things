package dm

import (
	"net/http"

	"github.com/i-Things/things/src/apisvr/internal/logic/open/dm"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenFirmwareUploadUriHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenFirmwareUploadUriReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dm.NewGenFirmwareUploadUriLogic(r.Context(), svcCtx)
		resp, err := l.GenFirmwareUploadUri(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
