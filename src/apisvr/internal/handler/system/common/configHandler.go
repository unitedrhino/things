package common

import (
	"net/http"

	"github.com/i-Things/things/shared/result"

	"github.com/i-Things/things/src/apisvr/internal/logic/system/common"
	"github.com/i-Things/things/src/apisvr/internal/svc"
)

func ConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewConfigLogic(r.Context(), svcCtx)
		resp, err := l.Config()
		result.Http(w, r, resp, err)
	}
}
