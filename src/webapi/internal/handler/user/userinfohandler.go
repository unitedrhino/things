package handler

import (
	"github.com/spf13/cast"
	"net/http"
	"gitee.com/godLei6/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
	"gitee.com/godLei6/things/src/webapi/internal/logic/user"
	"gitee.com/godLei6/things/src/webapi/internal/svc"
)

func UserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewUserInfoLogic(r.Context(), ctx)
		strUid := r.Header.Get(types.USER_UID)
		Uid := cast.ToInt64(strUid)
		resp, err := l.UserInfo(Uid)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
