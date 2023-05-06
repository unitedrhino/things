package authority

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"

	"github.com/i-Things/things/shared/result"

	"github.com/i-Things/things/src/apisvr/internal/logic/system/authority"
	"github.com/i-Things/things/src/apisvr/internal/svc"
)

func AuthorityApiIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthorityApiIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}
		l := authority.NewAuthorityApiIndexLogic(r.Context(), svcCtx)
		resp, err := l.AuthorityApiIndex(&req)
		result.Http(w, r, resp, err)
	}
}
