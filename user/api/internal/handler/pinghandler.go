package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"yl/user/api/internal/logic"
	"yl/user/api/internal/svc"
)

func pingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewPingLogic(r.Context(), ctx)
		err := l.Ping()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
