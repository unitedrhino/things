package handler

import (
	"net/http"

	"yl/user/api/internal/logic"
	"yl/user/api/internal/svc"
	"yl/user/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func register2Handler(ctx *svc.ServiceContext) http.HandlerFunc {
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
