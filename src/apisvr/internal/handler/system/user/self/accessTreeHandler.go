package self

import (
	"net/http"

	"github.com/i-Things/things/shared/result"

	"github.com/i-Things/things/src/apisvr/internal/logic/system/user/self"
	"github.com/i-Things/things/src/apisvr/internal/svc"
)

func AccessTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := self.NewAccessTreeLogic(r.Context(), svcCtx)
		resp, err := l.AccessTree()
		result.Http(w, r, resp, err)
	}
}
