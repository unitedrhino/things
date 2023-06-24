package middleware

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/userHeader"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/config"
	auth "github.com/i-Things/things/src/syssvr/client/auth"
	user "github.com/i-Things/things/src/syssvr/client/user"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type CheckTokenWareMiddleware struct {
	cfg     config.Config
	UserRpc user.User
	AuthRpc auth.Auth
}

func NewCheckTokenWareMiddleware(cfg config.Config, UserRpc user.User, AuthRpc auth.Auth) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{cfg: cfg, UserRpc: UserRpc, AuthRpc: AuthRpc}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("%s.Lifecycle.Before", utils.FuncName())

		var userCtx *userHeader.UserCtx

		isOpen, userCtx, err := m.OpenAuth(w, r)
		if isOpen { //如果是开放请求
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.OpenAuth error=%s", utils.FuncName(), err)
				http.Error(w, "开放请求失败："+err.Error(), http.StatusUnauthorized)
				return
			}
		} else { //如果是用户请求
			//校验 Jwt Token
			userCtx, err = m.UserAuth(w, r)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.UserAuth error=%s", utils.FuncName(), err)
				http.Error(w, "用户请求失败："+err.Error(), http.StatusUnauthorized)
				return
			}

			//校验 Casbin Rule
			_, err = m.AuthRpc.AuthApiCheck(r.Context(), &user.CheckAuthReq{
				RoleID: userCtx.Role,
				Path:   r.URL.Path,
				Method: utils.MethodToNum(r.Method),
			})
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
				http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
				return
			}
		}

		//注入 用户信息 到 ctx
		ctx2 := userHeader.SetUserCtx(r.Context(), userCtx)
		r = r.WithContext(ctx2)

		next(w, r)

		logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}

// 如果有开放认证的字段才进行认证
func (m *CheckTokenWareMiddleware) OpenAuth(w http.ResponseWriter, r *http.Request) (bool, *userHeader.UserCtx, error) {
	var isOpen bool
	userName, password, ok := r.BasicAuth()
	if !ok {
		return isOpen, nil, nil
	} else {
		isOpen = true
	}

	strIP, _ := utils.GetIP(r)
	if !m.cfg.OpenAuth.Auth(userName, password, strIP) {
		return isOpen, nil, errors.Permissions.AddMsg("开放认证没通过")
	}

	return isOpen, &userHeader.UserCtx{
		IsOpen:    isOpen,
		UserID:    0,
		Role:      0,
		IsAllData: true,
		IP:        strIP,
		Os:        r.Header.Get("User-Agent"),
	}, nil
}

func (m *CheckTokenWareMiddleware) UserAuth(w http.ResponseWriter, r *http.Request) (*userHeader.UserCtx, error) {
	strIP, _ := utils.GetIP(r)

	strToken := r.Header.Get(userHeader.UserTokenKey)
	if strToken == "" {
		logx.WithContext(r.Context()).Errorf("%s.CheckTokenWare ip=%s not find token",
			utils.FuncName(), strIP)
		return nil, errors.NotLogin
	}

	resp, err := m.UserRpc.UserCheckToken(r.Context(), &user.UserCheckTokenReq{
		Ip:    strIP,
		Token: strToken,
	})
	if err != nil {
		er := errors.Fmt(err)
		logx.WithContext(r.Context()).Errorf("%s.CheckTokenWare ip=%s token=%s return=%s",
			utils.FuncName(), strIP, strToken, err)
		return nil, er
	}

	if resp.Token != "" {
		w.Header().Set("Access-Control-Expose-Headers", userHeader.UserSetTokenKey)
		w.Header().Set(userHeader.UserSetTokenKey, resp.Token)
	}
	logx.WithContext(r.Context()).Infof("%s.CheckTokenWare ip:%v in.token=%s checkResp:%v",
		utils.FuncName(), strIP, strToken, utils.Fmt(resp))

	return &userHeader.UserCtx{
		IsOpen:    false,
		UserID:    resp.UserID,
		Role:      resp.Role,
		IsAllData: resp.IsAllData == def.True,
		IP:        strIP,
		Os:        r.Header.Get("User-Agent"),
	}, nil
}
