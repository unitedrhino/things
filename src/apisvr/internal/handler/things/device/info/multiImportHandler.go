package info

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device/info"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"path"
)

func MultiImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeviceMultiImportReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		f, fh, err := r.FormFile("file")
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("请上传文件:"+err.Error()))
			return
		}
		defer f.Close()

		//判断和限制格式
		if ext := path.Ext(fh.Filename); ext != ".csv" {
			result.Http(w, r, nil, errors.Parameter.WithMsg("仅支持csv文件"))
			return
		}

		//判断和限制大小
		if KB := float64(fh.Size) / 1024; KB > 700 { //byte->KB，限制最大700KB
			result.Http(w, r, nil, errors.Parameter.WithMsg("文件不能超过700KB"))
			return
		}

		fb := make([]byte, fh.Size)
		if _, err = f.Read(fb); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("解析文件字节失败:"+err.Error()))
			return
		}
		req.File = fb

		l := info.NewMultiImportLogic(r.Context(), svcCtx)
		resp, err := l.MultiImport(&req)
		result.Http(w, r, resp, err)
	}
}
