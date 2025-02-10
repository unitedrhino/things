package info

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device/info"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
)

// 创建绑定token
func BindTokenCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := info.NewBindTokenCreateLogic(r.Context(), svcCtx)
		resp, err := l.BindTokenCreate()
		result.Http(w, r, resp, err)
	}
}
