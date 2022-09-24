package frontend

import (
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"io"
	"net/http"
)

func FrontendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	dir := http.Dir(svcCtx.Config.FrontDir)
	rr, err := dir.Open("index.html")
	if err != nil {
		return func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusNotFound)
			writer.Write([]byte("404"))
			return
		}
	}

	file, err := io.ReadAll(rr)
	fileServer := http.FileServer(dir)

	return func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path

		if _, err := dir.Open(upath); err != nil {
			w.WriteHeader(http.StatusOK)
			w.Write(file)
			return
		}
		fileServer.ServeHTTP(w, r)
	}
}
