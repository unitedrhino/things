package middleware

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"strconv"
	"yl/shared/errors"
	"yl/shared/utils"
	"yl/src/user/rpc/user"
	"yl/src/user/rpc/userclient"
)

type CheckTokenMiddleware struct {
	UserRpc           userclient.User
}

func NewCheckTokenMiddleware(UserRpc           userclient.User) *CheckTokenMiddleware {
	return &CheckTokenMiddleware{UserRpc: UserRpc}
}

func (m *CheckTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strIP,_ := utils.GetIP(r)
		r.Header.Set("yl-ip",strIP)
		strToken := r.Header.Get("yl-token")
		resp,err:=m.UserRpc.CheckToken(r.Context(),&user.CheckTokenReq{
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
			w.Header().Set("yl-set-token",resp.Token)
		}
		strUid:= strconv.FormatInt(resp.Uid,10)
		r.Header.Set("yl-uid",strUid)
		logx.WithContext(r.Context()).Infof("CheckToken|ip=%s|uid=%s|token=%s|newToken=%s",
			strIP,strUid,strToken,resp.Token)
		next(w, r)
	}
}
