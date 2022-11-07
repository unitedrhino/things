package proxy

import (
	"fmt"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"io"
	"net/http"
)

func Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	dir := http.Dir(svcCtx.Config.Proxy.FrontDir)
	f, err := dir.Open(svcCtx.Config.Proxy.FrontDefaultPage)
	if err != nil {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte("404"))
			return
		}
	}

	indexFile, err := io.ReadAll(f)
	fileServer := http.FileServer(dir)

	return func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path
		f, err := dir.Open(upath)
		if err != nil {
			fmt.Println(f)
			w.WriteHeader(http.StatusOK)
			w.Write(indexFile)
			return
		} else {
			info, err := f.Stat()
			if err != nil || info.Mode().IsDir() {
				w.WriteHeader(http.StatusOK)
				w.Write(indexFile)
				return
			}
		}
		fileServer.ServeHTTP(w, r)
	}
}
