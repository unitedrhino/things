package info

import (
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device/info"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"github.com/spf13/cast"
	"net/http"
)

// 导入批量更新设备
func MultiUpdateImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			limitCnt = 5000 //限制表格数据条数
		)

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

		l := info.NewMultiUpdateImportLogic(r.Context(), svcCtx)
		resp, err := l.MultiUpdateImport(rows)
		result.Http(w, r, resp, err)
	}
}
