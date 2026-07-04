package securecontrol

import (
	"net/http"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	securecontrollogic "gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device/securecontrol"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// KeyGrantHandler 处理 App 安全控制 key-grant 请求。
func KeyGrantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		var req types.DeviceSecureControlKeyGrantReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := securecontrollogic.NewKeyGrantLogic(r.Context(), svcCtx)
		resp, err := l.KeyGrant(&req)
		result.Http(w, r, resp, err)
	}
}
