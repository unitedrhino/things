package hooks

import (
	"encoding/json"
	"gitee.com/i-Things/share/result"
	"github.com/i-Things/things/service/vidsvr/internal/logic/zlmedia/hooks"
	"github.com/i-Things/things/service/vidsvr/internal/svc"
	"github.com/i-Things/things/service/vidsvr/internal/types"
	"io"
	"net/http"
)

func OnStreamNoneReaderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HooksApiStreamNoneReaderReq
		//if err := httpx.Parse(r, &req); err != nil {
		//	result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
		//	return
		//}
		bodyByte, err := io.ReadAll(r.Body)
		json.Unmarshal(bodyByte, &req)
		l := hooks.NewOnStreamNoneReaderLogic(r.Context(), svcCtx)
		resp, err := l.OnStreamNoneReader(&req)
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
