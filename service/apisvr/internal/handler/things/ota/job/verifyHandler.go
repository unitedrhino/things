package job

import (
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/result"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/ota/job"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func VerifyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OtaFirmwareVerifyReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := job.NewVerifyLogic(r.Context(), svcCtx)
		resp, err := l.Verify(&req)
		result.Http(w, r, resp, err)
	}
}
