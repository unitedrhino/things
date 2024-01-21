package proxy

import (
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	dir := http.Dir(svcCtx.Config.Proxy.FileProxy.FrontDir)
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
			defaultHandle(svcCtx, upath, w, r, false)
			return
		} else {
			info, err := f.Stat()
			if err != nil || info.Mode().IsDir() {
				defaultHandle(svcCtx, upath, w, r, info.Mode().IsDir())
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
	if conf.DeletePrefix {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, conf.Router)
	}
	proxy.ServeHTTP(w, r)
}
func defaultHandle(svcCtx *svc.ServiceContext, upath string, w http.ResponseWriter, r *http.Request, isDir bool) {
	dir := http.Dir(svcCtx.Config.Proxy.FileProxy.FrontDir)
	var (
		f   http.File
		err error
	)
	switch upath {
	case "/favicon.ico":
		f, err = dir.Open(fmt.Sprintf("%sfavicon.ico", svcCtx.Config.Proxy.FileProxy.CoreDir))
	case "/":
		http.Redirect(w, r, svcCtx.Config.Proxy.FileProxy.CoreDir, http.StatusMovedPermanently)
		return
	default:
		if strings.HasSuffix(upath, "/") {
			f, err = dir.Open(upath + "index.html")
		} else if isDir {
			f, err = dir.Open(upath + "/index.html")
		}
	}
	if err != nil || f == nil {
		http.NotFound(w, r)
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
