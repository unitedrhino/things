package pair

import (
	"net/http"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/result"
	pairlogic "gitee.com/unitedrhino/things/service/apisvr/internal/logic/things/device/pair"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GrantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DevicePairGrantReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := pairlogic.NewGrantLogic(r.Context(), svcCtx)
		resp, err := l.Grant(&req)
		result.Http(w, r, resp, err)
	}
}
