package collect

import (
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/result"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/user/device/collect"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func MultiCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserCollectDeviceSave
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := collect.NewMultiCreateLogic(r.Context(), svcCtx)
		err := l.MultiCreate(&req)
		result.Http(w, r, nil, err)
	}
}
