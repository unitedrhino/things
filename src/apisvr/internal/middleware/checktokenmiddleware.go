package middleware

import (
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/config"
	"github.com/i-Things/things/src/apisvr/internal/domain/userHeader"
	user "github.com/i-Things/things/src/syssvr/client/user"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type CheckTokenMiddleware struct {
	UserRpc user.User
	c       config.Config
}

func NewCheckTokenMiddleware(c config.Config, UserRpc user.User) *CheckTokenMiddleware {
	return &CheckTokenMiddleware{UserRpc: UserRpc, c: c}
}

func (m *CheckTokenMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err, isOpen := m.OpenAuth(w, r)
		if isOpen { //如果是开放请求
			if err == nil {
				next(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
			return
		}
		userCtx, err := m.UserAuth(w, r)
		if err == nil {
			next(w, r.WithContext(userHeader.SetUserCtx(r.Context(), userCtx)))
			return
		}
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}

// 如果有开放认证的字段才进行认证
func (m *CheckTokenMiddleware) OpenAuth(w http.ResponseWriter, r *http.Request) (error, bool) {
	userName, password, ok := r.BasicAuth()
	if !ok {
		return nil, false
	}
	strIP, _ := utils.GetIP(r)
	if !m.c.OpenAuth.Auth(userName, password, strIP) {
		return errors.Permissions.AddMsg("开放认证没通过"), true
	}
	return nil, true
}

func (m *CheckTokenMiddleware) UserAuth(w http.ResponseWriter, r *http.Request) (*userHeader.UserCtx, error) {
	strIP, _ := utils.GetIP(r)
	strToken := r.Header.Get(userHeader.UserToken)
	if strToken == "" {
		logx.WithContext(r.Context()).Errorf("%s.CheckToken ip=%s not find token",
			utils.FuncName(), strIP)
		return nil, errors.NotLogin
	}
	resp, err := m.UserRpc.CheckToken(r.Context(), &user.CheckTokenReq{
		Ip:    strIP,
		Token: strToken,
	})
	if err != nil {
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("%s.CheckToken ip=%s token=%s return=%s",
			utils.FuncName(), strIP, strToken, err)
		return nil, er
	}
	if resp.Token != "" {
		w.Header().Set("Access-Control-Expose-Headers", userHeader.UserSetToken)
		w.Header().Set(userHeader.UserSetToken, resp.Token)
	}
	logx.WithContext(r.Context()).Infof("%s.CheckToken ip:%v in.token=%s checkResp:%v",
		utils.FuncName(), strIP, strToken, utils.Fmt(resp))
	return &userHeader.UserCtx{
		Uid:  resp.Uid,
		IP:   strIP,
		Role: resp.Role,
	}, nil
}
