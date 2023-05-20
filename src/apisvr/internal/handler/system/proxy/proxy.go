package proxy

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	dir := http.Dir(svcCtx.Config.Proxy.FileProxy[0].FrontDir)
	_, err := dir.Open(svcCtx.Config.Proxy.FileProxy[0].FrontDefaultPage)
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
		for _, v := range svcCtx.Config.Proxy.StaticProxy {
			if strings.HasPrefix(upath, v.Router) {
				staticProxy(svcCtx, v, w, r)
				return
			}
		}
		f, err := dir.Open(upath)
		if err != nil {
			defaultHandle(svcCtx, upath, w, r)
			return
		} else {
			info, err := f.Stat()
			if err != nil || info.Mode().IsDir() {
				defaultHandle(svcCtx, upath, w, r)
				return
			}
		}
		fileServer.ServeHTTP(w, r)
	}
}

func staticProxy(svcCtx *svc.ServiceContext, conf *conf.StaticProxyConf, w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse(conf.Dest)
	if err != nil {
		//defaultHandle(svcCtx, w, r)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	r.Host = remote.Host
	proxy.ServeHTTP(w, r)
}
func defaultHandle(svcCtx *svc.ServiceContext, upath string, w http.ResponseWriter, r *http.Request) {
	dir := http.Dir(svcCtx.Config.Proxy.FileProxy[0].FrontDir)
	f, err := dir.Open(svcCtx.Config.Proxy.FileProxy[0].FrontDefaultPage)
	if upath == "/favicon.ico" {
		f, err = dir.Open("front/iThingsCore/favicon.ico")
	}
	if err != nil {
		return
	}
	indexFile, err := io.ReadAll(f)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(indexFile)
	return
}
