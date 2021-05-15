package handler

import (
	"net/http"

	"yl/src/webapi/internal/logic/verify"
	"yl/src/webapi/internal/svc"
	"yl/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func CaptchaHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCaptchaReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewCaptchaLogic(r.Context(), ctx)
		resp, err := l.Captcha(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
