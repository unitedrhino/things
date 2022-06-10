package oss

import (
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func UploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		if true {
			remote, err := url.Parse("http://127.0.0.1:8567")
			if err != nil {
				panic(err)
			}
			proxy := httputil.NewSingleHostReverseProxy(remote)

			originalDirector := proxy.Director
			proxy.Director = func(hReq *http.Request) {
				originalDirector(hReq)
				hReq.URL.Path = "/oss/upload"
				hReq.Header.Set("I-Things-Business", req.Business)
			}
			proxy.ServeHTTP(w, r)
			return
		}
	}
}
