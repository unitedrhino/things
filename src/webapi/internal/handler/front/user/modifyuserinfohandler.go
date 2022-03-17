package user

import (
	"github.com/i-Things/things/src/webapi/internal/logic/front/user"
	"github.com/spf13/cast"
	"net/http"

	"github.com/i-Things/things/src/webapi/internal/svc"
	"github.com/i-Things/things/src/webapi/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ModifyUserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModifyUserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		strUid := r.Header.Get(types.USER_UID)
		Uid := cast.ToInt64(strUid)
		l := user.NewModifyUserInfoLogic(r.Context(), ctx)
		err := l.ModifyUserInfo(req, Uid)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
