package user

import (
	"net/http"

	"github.com/go-things/things/src/webapi/internal/logic/backgrand/user"
	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func UserInfosHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserInfosReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewUserInfosLogic(r.Context(), ctx)
		resp, err := l.UserInfos(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
