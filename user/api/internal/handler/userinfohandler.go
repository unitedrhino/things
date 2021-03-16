package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"yl/user/api/internal/logic"
	"yl/user/api/internal/svc"
)

func userInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewUserInfoLogic(r.Context(), ctx)
		resp, err := l.UserInfo()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
