package self

import (
	"net/http"

	"github.com/i-Things/things/shared/result"

	"github.com/i-Things/things/src/apisvr/internal/logic/system/user/self"
	"github.com/i-Things/things/src/apisvr/internal/svc"
)

func ResourceReadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := self.NewResourceReadLogic(r.Context(), svcCtx)
		resp, err := l.ResourceRead()
		result.Http(w, r, resp, err)
	}
}
