package handler

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
	"yl/src/webapi/api/test/internal/logic/user"
	"yl/src/webapi/api/test/internal/svc"
)

func PingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
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
