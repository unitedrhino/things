package auth

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"gitee.com/i-Things/things/service/apisvr/internal/logic/things/device/auth"
	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func AccessHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeviceAuthAccessReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.AddMsg(err.Error()))
			return
		}

		l := auth.NewAccessLogic(r.Context(), svcCtx)
		err := l.Access(&req)
		result.Http(w, r, nil, err)
	}
}
