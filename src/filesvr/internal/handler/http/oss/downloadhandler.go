package oss

import (
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/minio/minio-go/v7"
	"net/http"
	"time"

	"github.com/i-Things/things/src/filesvr/internal/svc"
	"github.com/i-Things/things/src/filesvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownLoadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		err := func() error {
			obj, err := svcCtx.OSS.GetObject(r.Context(), "mymusic", req.Sign, minio.GetObjectOptions{})
			fmt.Println(obj, err)
			if err != nil {
				return errors.Fmt(err)
			}
			w.Header().Set("Content-Disposition", "attachment; filename="+req.Sign)
			http.ServeContent(w, r, req.Sign, time.Now(), obj)
			return nil
		}()
		if err != nil {
			httpx.Error(w, err)
		}
	}
}
