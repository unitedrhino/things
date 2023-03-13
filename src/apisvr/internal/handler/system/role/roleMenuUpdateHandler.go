package role

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/result"
	"github.com/i-Things/things/src/apisvr/internal/logic/system/role"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func RoleMenuUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleMenuUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.AddMsg(err.Error()))
			return
		}

		l := role.NewRoleMenuUpdateLogic(r.Context(), svcCtx)
		err := l.RoleMenuUpdate(&req)
		result.Http(w, r, nil, err)
	}
}
