package handler

import (
	"net/http"
	"yl/shared/utils"
	"yl/user/api/internal/logic/user"
	"yl/user/api/internal/svc"
	"yl/user/api/internal/types"
	"yl/user/common"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func Register2Handler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Register2Req
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		jwt := r.Header.Get("authorization")
		token,err:=utils.ParseToken(jwt,ctx.Config.Rej.AccessSecret)
		if err != nil {
			httpx.Error(w, err)
			return
		}
		if token.Uid != req.Uid {
			httpx.Error(w,common.ErrorUidNotCompare)
			return
		}
		l := logic.NewRegister2Logic(r.Context(), ctx)
		err = l.Register2(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
