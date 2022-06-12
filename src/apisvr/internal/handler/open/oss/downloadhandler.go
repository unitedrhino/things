package oss

import (
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"time"
)

func DownLoadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		err := func() error {
			token, err := devices.ParseToken(req.Sign, svcCtx.Config.OSS.AccessSecret)
			if err != nil {
				return errors.Fmt(err)
			}
			stat, err := svcCtx.OSS.StatObject(r.Context(), token.Bucket, token.Dir, minio.StatObjectOptions{})
			if err != nil {
				return errors.Fmt(err)
			}
			fileName := stat.UserMetadata["Filename"]
			obj, err := svcCtx.OSS.GetObject(r.Context(), token.Bucket, token.Dir, minio.GetObjectOptions{})
			if err != nil {
				return errors.Fmt(err)
			}
			w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
			http.ServeContent(w, r, fileName, time.Now(), obj)
			return nil
		}()
		if err != nil {
			httpx.Error(w, err)
		}

	}
}
