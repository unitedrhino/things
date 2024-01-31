package gbsip

import (
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/result"
	"github.com/i-Things/things/service/apisvr/internal/logic/things/vidmgr/gbsip"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func DeletedevHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.VidmgrSipDeleteDevReq
		if err := httpx.Parse(r, &req); err != nil {
			result.Http(w, r, nil, errors.Parameter.WithMsg("入参不正确:"+err.Error()))
			return
		}

		l := gbsip.NewDeletedevLogic(r.Context(), svcCtx)
		err := l.Deletedev(&req)
		result.Http(w, r, nil, err)
	}
}
