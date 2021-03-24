package middleware

import (
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"yl/shared/errors"
	"yl/src/user/rpc/user"
	"yl/src/user/rpc/userclient"
)

type UsercheckMiddleware struct {
	UserRpc           userclient.User
}

func NewUsercheckMiddleware(UserRpc           userclient.User) *UsercheckMiddleware {
	return &UsercheckMiddleware{UserRpc: UserRpc}
}

func (m *UsercheckMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		resp,err:=m.UserRpc.CheckToken(r.Context(),&user.CheckTokenReq{})
		if err != nil {
			er := errors.Fmt(err)
			logx.WithContext(r.Context()).Errorf("get error=%#v|fmt error=%#v",err,er)
		}
		logx.WithContext(r.Context()).Infof("resp=%#v",resp)
		// Passthrough to next handler if need
		next(w, r)
	}
}
