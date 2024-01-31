package interact

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/device/interact"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func SendActionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeviceInteractSendActionReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.AddMsg(err.Error()))
			return
		}

		l := interact.NewSendActionLogic(r.Context(), svcCtx)
		resp, err := l.SendAction(&req)
		result.Http(w, r, resp, err)
	}
}
