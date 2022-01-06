package dm

import (
	"github.com/go-things/things/src/webapi/internal/logic/backgrand/dm"
	"net/http"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func LoginAuthHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginAuthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dm.NewLoginAuthLogic(r.Context(), ctx)
		err := l.LoginAuth(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
