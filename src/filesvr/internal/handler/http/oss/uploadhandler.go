package oss

import (
	"net/http"

	"github.com/i-Things/things/src/filesvr/internal/logic/http/oss"
	"github.com/i-Things/things/src/filesvr/internal/svc"
	"github.com/i-Things/things/src/filesvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := oss.NewUploadLogic(r.Context(), svcCtx)
		err := l.Upload(&req, r)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
