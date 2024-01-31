package hooks

import (
	"encoding/json"
	"gitee.com/i-Things/core/shared/result"
	"github.com/i-Things/things/src/vidsvr/internal/logic/zlmedia/hooks"
	"github.com/i-Things/things/src/vidsvr/internal/svc"
	"github.com/i-Things/things/src/vidsvr/internal/types"
	"io"
	"net/http"
)

func OnShellLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HooksApiShellLoginReq
		//if err := httpx.Parse(r, &req); err != nil {
		//	result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
		//	return
		//}
		bodyByte, err := io.ReadAll(r.Body)
		json.Unmarshal(bodyByte, &req)
		l := hooks.NewOnShellLoginLogic(r.Context(), svcCtx)
		resp, err := l.OnShellLogin(&req)
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
