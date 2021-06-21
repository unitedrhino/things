package handler

import (
	"github.com/spf13/cast"
	"net/http"

	"gitee.com/godLei6/things/src/webapi/internal/logic/user"
	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
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
		l := logic.NewModifyUserInfoLogic(r.Context(), ctx)
		err := l.ModifyUserInfo(req, Uid)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
