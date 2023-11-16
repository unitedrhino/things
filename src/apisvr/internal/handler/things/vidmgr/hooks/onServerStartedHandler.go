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

func OnServerStartedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HooksApiServerStartedReq
		//if err := httpx.Parse(r, &req); err != nil {
		//	result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
		//	return
		//}

		//httpx.Parse在解析XXXX.XXXX的json数据格式会出错，go-zero框架代码问题。
		//采用如下代码
		bodyByte, err := io.ReadAll(r.Body)
		json.Unmarshal(bodyByte, &req)
		//fmt.Println("[--Debug--]", string(bodyByte))

		l := hooks.NewOnServerStartedLogic(r.Context(), svcCtx)
		resp, err := l.OnServerStarted(&req)
		result.HttpWithoutWrap(w, r, resp, err)
	}
}
