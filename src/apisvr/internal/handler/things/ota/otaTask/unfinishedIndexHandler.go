package otaTask

import (
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/result"
	"github.com/i-Things/things/src/apisvr/internal/logic/things/ota/otaTask"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func UnfinishedIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OTAUnfinishedTaskByDeviceIndexReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := otaTask.NewUnfinishedIndexLogic(r.Context(), svcCtx)
		resp, err := l.UnfinishedIndex(&req)
		result.Http(w, r, resp, err)
	}
}
