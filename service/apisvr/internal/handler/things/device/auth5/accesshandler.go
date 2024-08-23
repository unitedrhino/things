package auth5

import (
	"gitee.com/i-Things/share/result"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/device/auth5"
	"net/http"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AccessHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeviceAuth5AccessReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := auth5.NewAccessLogic(r.Context(), svcCtx)
		resp, err := l.Access(&req)
		l.Infof("%s req=%v resp=%v err=%v", utils.FuncName(), utils.Fmt(req), utils.Fmt(resp), err)
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
