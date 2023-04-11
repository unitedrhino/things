package info

import (
	"bytes"
	"encoding/csv"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device/info"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/rest/httpx"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
	"path"
)

func MultiImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				utils.HandleThrow(r.Context(), p)
				result.Http(w, r, nil, errors.Panic)
				return
			}
		}()

		var (
			req      types.DeviceMultiImportReq
			limitCnt = 1000              //限制表格数据条数
			limitKB  = float64(5 * 1024) //限制表格文件大小（5M）
		)

		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		f, fh, err := r.FormFile("file")
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("请上传csv文件").AddDetail(err.Error()))
			return
		}
		defer f.Close()

		//判断和限制格式
		if ext := path.Ext(fh.Filename); ext != ".csv" {
			result.Http(w, r, nil, errors.Parameter.WithMsg("仅支持csv文件"))
			return
		}

		//判断和限制大小
		if fileKB := float64(fh.Size) / 1024; fileKB > limitKB { //byte->KB
			result.Http(w, r, nil, errors.Parameter.WithMsgf("文件不能超过%.2fM", limitKB/1024))
			return
		}

		fb := make([]byte, fh.Size)
		if _, err = f.Read(fb); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取表格字节失败:"+err.Error()))
			return
		}

		reader := csv.NewReader(transform.NewReader(bytes.NewReader(fb), simplifiedchinese.GB18030.NewDecoder()))
		rows, err := reader.ReadAll()
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取表格内容失败:"+err.Error()))
			return
		}
		if len(rows)-1 > limitCnt {
			result.Http(w, r, nil, errors.Parameter.WithMsgf("最多只能导入%s条数据", cast.ToString(limitCnt)))
			return
		}

		req.File = fb
		l := info.NewMultiImportLogic(r.Context(), svcCtx)
		resp, err := l.MultiImport(&req, rows)
		result.Http(w, r, resp, err)
	}
}
