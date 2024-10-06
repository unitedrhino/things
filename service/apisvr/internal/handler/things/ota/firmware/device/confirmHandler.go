package device

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"gitee.com/i-Things/things/service/apisvr/internal/logic/things/ota/firmware/device"
	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func ConfirmHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OtaFirmwareDeviceConfirmReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := device.NewConfirmLogic(r.Context(), svcCtx)
		err := l.Confirm(&req)
		result.Http(w, r, nil, err)
	}
}
