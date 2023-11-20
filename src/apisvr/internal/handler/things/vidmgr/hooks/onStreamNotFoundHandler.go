package hooks

import (
	"encoding/json"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/vidmgr/hooks"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"io"
	"net/http"
)

func OnStreamNotFoundHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HooksApiStreamNotFoundReq
		//if err := httpx.Parse(r, &req); err != nil {
		//	result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
		//	return
		//}
		bodyByte, err := io.ReadAll(r.Body)
		json.Unmarshal(bodyByte, &req)
		l := hooks.NewOnStreamNotFoundLogic(r.Context(), svcCtx)
		resp, err := l.OnStreamNotFound(&req)
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
