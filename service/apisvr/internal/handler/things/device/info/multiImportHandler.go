package info

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/device/info"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
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
			limitCnt = 1000 //限制表格数据条数
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
		rows, err := utils.ReadExcel(f, fh.Filename)
		if err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("读取表格内容失败:"+err.Error()))
			return
		}

		if len(rows)-1 > limitCnt {
			result.Http(w, r, nil, errors.Parameter.WithMsgf("最多只能导入%s条数据", cast.ToString(limitCnt)))
			return
		}
		l := info.NewMultiImportLogic(r.Context(), svcCtx)
		resp, err := l.MultiImport(&req, rows)
		result.Http(w, r, resp, err)
	}
}
