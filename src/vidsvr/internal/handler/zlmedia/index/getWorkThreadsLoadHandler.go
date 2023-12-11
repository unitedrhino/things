package indexapi

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/src/vidsvr/internal/logic/zlmedia/index"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func GetWorkThreadsLoadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IndexApiReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := index.NewGetWorkThreadsLoadLogic(r.Context(), svcCtx)
		resp, err := l.GetWorkThreadsLoad(&req)
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
