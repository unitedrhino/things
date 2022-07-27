package device

import (
	"net/http"

	"github.com/i-Things/things/shared/result"

	"github.com/i-Things/things/src/apisvr/internal/logic/things/device"
	"github.com/i-Things/things/src/apisvr/internal/svc"
)

func TestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := device.NewTestLogic(r.Context(), svcCtx)
		err := l.Test()
		result.Http(w, r, nil, err)
	}
}
