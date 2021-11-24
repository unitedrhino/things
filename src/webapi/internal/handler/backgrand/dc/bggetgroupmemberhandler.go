package dc

import (
	"net/http"

	"gitee.com/godLei6/things/src/webapi/internal/logic/backgrand/dc"
	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func BgGetGroupMemberHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGroupMemberReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := dc.NewBgGetGroupMemberLogic(r.Context(), ctx)
		resp, err := l.BgGetGroupMember(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
