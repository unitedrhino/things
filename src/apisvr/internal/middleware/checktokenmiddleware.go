package middleware

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	"github.com/i-Things/things/src/usersvr/pb/user"
	"github.com/i-Things/things/src/usersvr/userclient"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type CheckTokenMiddleware struct {
	UserRpc userclient.User
}

func NewCheckTokenMiddleware(UserRpc userclient.User) *CheckTokenMiddleware {
	return &CheckTokenMiddleware{UserRpc: UserRpc}
}

func (m *CheckTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strIP, _ := utils.GetIP(r)
		strToken := r.Header.Get(userHeader.UserToken)
		if strToken == "" {
			logx.WithContext(r.Context()).Errorf("%s|CheckToken|ip=%s|not find token",
				utils.FuncName(), strIP)
			http.Error(w, errors.TokenMalformed.Error(), http.StatusUnauthorized)
			return
		}
		resp, err := m.UserRpc.CheckToken(r.Context(), &user.CheckTokenReq{
			Ip:    strIP,
			Token: strToken,
		})
		if err != nil {
			er := errors.Fmt(err)
			logx.WithContext(r.Context()).Errorf("%s|CheckToken|ip=%s|token=%s|return=%s",
				utils.FuncName(), strIP, strToken, err)
			http.Error(w, er.Error(), http.StatusUnauthorized)
			return
		}
		if resp.Token != "" {
			w.Header().Set(userHeader.UserSetToken, resp.Token)
		}
		logx.WithContext(r.Context()).Infof("CheckToken|ip=%s|uid=%s|token=%s|newToken=%s",
			strIP, resp.Uid, strToken, resp.Token)
		next(w, r.WithContext(userHeader.SetUserCtx(r.Context(), resp.Uid, strIP)))
	}
}
