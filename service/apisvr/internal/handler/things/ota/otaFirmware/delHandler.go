package otaFirmware

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/ota/otaFirmware"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func DelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FirmwareDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}
		l := otaFirmware.NewDelLogic(r.Context(), svcCtx)
		err := l.Del(&req)
		result.Http(w, r, nil, err)
	}
}
