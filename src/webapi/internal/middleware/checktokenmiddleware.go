package middleware

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"strconv"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/user"
	"yl/src/user/userclient"
	"yl/src/webapi/internal/types"
)

type CheckTokenMiddleware struct {
	UserRpc userclient.User
}

func NewCheckTokenMiddleware(UserRpc userclient.User) *CheckTokenMiddleware {
	return &CheckTokenMiddleware{UserRpc: UserRpc}
}

func (m *CheckTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strIP,_ := utils.GetIP(r)
		r.Header.Set(types.USER_IP,strIP)
		strToken := r.Header.Get(types.USER_TOKEN)
		if strToken == "" {
			logx.WithContext(r.Context()).Errorf("%s|CheckToken|ip=%s|not find token",
				utils.FuncName(),strIP)
			httpx.Error(w, errors.TokenMalformed)
			return
		}
		resp,err := m.UserRpc.CheckToken(r.Context(),&user.CheckTokenReq{
			Ip: strIP,
			Token: strToken,
		})
		if err != nil {
			er := errors.Fmt(err)
			logx.WithContext(r.Context()).Errorf("%s|CheckToken|ip=%s|token=%s|return=%s",
				utils.FuncName(),strIP,strToken,err)
			httpx.Error(w, er.AddDetail("token检查未通过"))
			return
		}
		if resp.Token != "" {
			w.Header().Set(types.USER_SET_TOKEN,resp.Token)
		}
		strUid:= strconv.FormatInt(resp.Uid,10)
		r.Header.Set(types.USER_UID,strUid)
		logx.WithContext(r.Context()).Infof("CheckToken|ip=%s|uid=%s|token=%s|newToken=%s",
			strIP,strUid,strToken,resp.Token)
		next(w, r)
	}
}
