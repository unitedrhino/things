package authority

import (
	"net/http"

	"github.com/i-Things/things/shared/result"

	"github.com/i-Things/things/src/apisvr/internal/logic/system/authority"
	"github.com/i-Things/things/src/apisvr/internal/svc"
)

func AuthorityApiIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := authority.NewAuthorityApiIndexLogic(r.Context(), svcCtx)
		resp, err := l.AuthorityApiIndex()
		result.Http(w, r, resp, err)
	}
}
