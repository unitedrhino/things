package oss

import (
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func DownLoadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		remote, err := url.Parse("http://127.0.0.1:8567")
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(remote)

		originalDirector := proxy.Director
		proxy.Director = func(hReq *http.Request) {
			originalDirector(hReq)
			hReq.URL.Path = "/oss/download"
			hReq.Header.Set("I-Things-Business", req.Business)
		}
		proxy.ServeHTTP(w, r)
		return
	}
}
func test(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	//w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	//w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	//http.ServeContent(rw, c.Ctx.Request, "(文件名字)", time.Now(), bytes.NewReader(b.Bytes()))
}
