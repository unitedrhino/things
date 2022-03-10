package dc

import (
	"net/http"

	"github.com/go-things/things/src/webapi/internal/logic/open/dc"
	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func BgManageGroupMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ManageGroupMemberReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dc.NewBgManageGroupMemberLogic(r.Context(), svcCtx)
		resp, err := l.BgManageGroupMember(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
