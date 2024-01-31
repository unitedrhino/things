package hooks

import (
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/result"
	"github.com/i-Things/things/src/vidsvr/internal/logic/zlmedia/hooks"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func OnRecordMp4Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HooksApiRecordMp4Req
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := hooks.NewOnRecordMp4Logic(r.Context(), svcCtx)
		resp, err := l.OnRecordMp4(&req)
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
