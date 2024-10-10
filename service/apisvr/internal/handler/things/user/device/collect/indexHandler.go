package collect

import (
	"net/http"

	"gitee.com/unitedrhino/share/result"

	"gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/user/device/collect"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
)

func IndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := collect.NewIndexLogic(r.Context(), svcCtx)
		resp, err := l.Index()
		result.Http(w, r, resp, err)
	}
}
