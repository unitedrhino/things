package handler

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"yl/src/user/api/internal/logic/user"
	"yl/src/user/api/internal/svc"
	"yl/src/user/api/internal/types"
)

func Register2Handler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Register2Req
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		l := logic.NewRegister2Logic(r.Context(), ctx)
		err := l.Register2(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
