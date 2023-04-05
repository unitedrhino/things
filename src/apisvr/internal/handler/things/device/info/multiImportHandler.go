package info

import (
	"bytes"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/shared/utils/cast"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/device/info"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/xuri/excelize/v2"
	"github.com/zeromicro/go-zero/rest/httpx"
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

		demoCnt := int64(3)
		limitCnt := int64(1000)
		rowLimitCnt := demoCnt + limitCnt

		var req types.DeviceMultiImportReq
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
		if KB := float64(fh.Size) / 1024; KB > 700 { //byte->KB，限制最大700KB
			result.Http(w, r, nil, errors.Parameter.WithMsg("文件不能超过700KB"))
			return
		}

		fb := make([]byte, fh.Size)
		if _, err = f.Read(fb); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取表格字节失败:"+err.Error()))
			return
		}

		reader := bytes.NewReader(fb)
		csv, err := excelize.OpenReader(reader)
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取表格内容失败:"+err.Error()))
			return
		}
		rows, err := csv.Rows(csv.GetSheetName(0))
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取表格Sheet失败:"+err.Error()))
			return
		}

		for rowCnt := int64(0); rows.Next(); rowCnt++ {
			if rowCnt >= rowLimitCnt {
				result.Http(w, r, nil, errors.Parameter.WithMsgf("最多只能导入%s条数据", cast.ToString(limitCnt)))
				return
			}
		}

		req.File = fb
		l := info.NewMultiImportLogic(r.Context(), svcCtx)
		resp, err := l.MultiImport(&req, csv, demoCnt)
		result.Http(w, r, resp, err)
	}
}
