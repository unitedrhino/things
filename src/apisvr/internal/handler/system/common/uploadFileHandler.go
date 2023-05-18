package common

import (
	"net/http"

	"github.com/i-Things/things/shared/result"

	"github.com/i-Things/things/src/apisvr/internal/logic/system/common"
	"github.com/i-Things/things/src/apisvr/internal/svc"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := common.NewUploadFileLogic(r.Context(), svcCtx, r)
		resp, err := l.UploadFile()
		result.Http(w, r, resp, err)
	}
}
