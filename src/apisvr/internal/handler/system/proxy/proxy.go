package proxy

import (
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"io"
	"net/http"
)

func Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	dir := http.Dir(svcCtx.Config.Proxy.FrontDir)
	_, err := dir.Open(svcCtx.Config.Proxy.FrontDefaultPage)
	if err != nil { //没有前端代理模式
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte("404"))
			return
		}
	}
	fileServer := http.FileServer(dir)
	return func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path
		f, err := dir.Open(upath)
		if err != nil {
			defaultHandle(svcCtx, w, r)
			return
		} else {
			info, err := f.Stat()
			if err != nil || info.Mode().IsDir() {
				defaultHandle(svcCtx, w, r)
				return
			}
		}
		fileServer.ServeHTTP(w, r)
	}
}
func defaultHandle(svcCtx *svc.ServiceContext, w http.ResponseWriter, r *http.Request) {
	dir := http.Dir(svcCtx.Config.Proxy.FrontDir)
	f, err := dir.Open(svcCtx.Config.Proxy.FrontDefaultPage)
	if err != nil {
		return
	}
	indexFile, err := io.ReadAll(f)
	w.WriteHeader(http.StatusOK)
	w.Write(indexFile)
	return
}
