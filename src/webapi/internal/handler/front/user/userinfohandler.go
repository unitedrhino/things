package user

import (
	"gitee.com/godLei6/things/src/webapi/internal/logic/front/user"
	"net/http"

	"gitee.com/godLei6/things/src/webapi/internal/svc"
	"github.com/tal-tech/go-zero/rest/httpx"
)

func UserInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := user.NewUserInfoLogic(r.Context(), ctx)
		resp, err := l.UserInfo()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
