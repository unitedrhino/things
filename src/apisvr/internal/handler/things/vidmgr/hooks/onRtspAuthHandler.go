package hooks

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/vidmgr/hooks"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func OnRtspAuthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HooksApiPublishReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := hooks.NewOnRtspAuthLogic(r.Context(), svcCtx)
		resp, err := l.OnRtspAuth(&req)
		result.HooksApiHttp(w, r, resp, err)
	}
}
