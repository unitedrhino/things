package remoteConfig

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/product/remoteConfig"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func PushAllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ProductRemoteConfigPushAllReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := remoteConfig.NewPushAllLogic(r.Context(), svcCtx)
		err := l.PushAll(&req)
		result.Http(w, r, nil, err)
	}
}
