package middleware

import (
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/config"
	role "github.com/i-Things/things/src/syssvr/client/rolemanage"
	user "github.com/i-Things/things/src/syssvr/client/usermanage"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type CheckTokenWareMiddleware struct {
	cfg     config.Config
	UserRpc user.UserManage
	AuthRpc role.RoleManage
}

func NewCheckTokenWareMiddleware(cfg config.Config, UserRpc user.UserManage, AuthRpc role.RoleManage) *CheckTokenWareMiddleware {
	return &CheckTokenWareMiddleware{cfg: cfg, UserRpc: UserRpc, AuthRpc: AuthRpc}
}

func (m *CheckTokenWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("%s.Lifecycle.Before", utils.FuncName())

		var userCtx *ctxs.UserCtx

		isOpen, userCtx, err := m.OpenAuth(w, r)
		if isOpen { //如果是开放请求
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.OpenAuth error=%s", utils.FuncName(), err)
				http.Error(w, "开放请求失败："+err.Error(), http.StatusUnauthorized)
				return
			}
			//注入 用户信息 到 ctx
			ctx2 := ctxs.SetUserCtx(r.Context(), userCtx)
			r = r.WithContext(ctx2)
		} else { //如果是用户请求
			//校验 Jwt Token
			userCtx, err = m.UserAuth(w, r)
			if err != nil {
				logx.WithContext(r.Context()).Errorf("%s.UserAuth error=%s", utils.FuncName(), err)
				http.Error(w, "用户请求失败："+err.Error(), http.StatusUnauthorized)
				return
			}
			//注入 用户信息 到 ctx
			ctx2 := ctxs.SetUserCtx(r.Context(), userCtx)
			r = r.WithContext(ctx2)
			////校验 Casbin Rule
			//_, err = m.AuthRpc.RoleApiAuth(r.Context(), &user.RoleApiAuthReq{
			//	RoleID: userCtx.RoleID,
			//	Path:   r.URL.Path,
			//	Method: r.Method,
			//})
			//if err != nil {
			//	logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
			//	http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
			//	return
			//}
		}

		next(w, r)

		logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}

// 如果有开放认证的字段才进行认证
func (m *CheckTokenWareMiddleware) OpenAuth(w http.ResponseWriter, r *http.Request) (bool, *ctxs.UserCtx, error) {
	var isOpen bool
	userName, password, ok := r.BasicAuth()
	if !ok {
		return isOpen, nil, nil
	} else {
		isOpen = true
	}

	strIP, _ := utils.GetIP(r)
	if !utils.Auth(m.cfg.OpenAuth, userName, password, strIP) {
		return isOpen, nil, errors.Permissions.AddMsg("开放认证没通过")
	}

	return isOpen, &ctxs.UserCtx{
		IsOpen:    isOpen,
		UserID:    0,
		RoleID:    0,
		IsAllData: true,
		IP:        strIP,
		Os:        r.Header.Get("User-Agent"),
	}, nil
}

func (m *CheckTokenWareMiddleware) UserAuth(w http.ResponseWriter, r *http.Request) (*ctxs.UserCtx, error) {
	strIP, _ := utils.GetIP(r)

	strToken := r.Header.Get(ctxs.UserTokenKey)
	if strToken == "" {
		logx.WithContext(r.Context()).Errorf("%s.CheckTokenWare ip=%s not find token",
			utils.FuncName(), strIP)
		return nil, errors.NotLogin
	}
	strRoleID := r.Header.Get(ctxs.UserRoleKey)
	roleID := cast.ToInt64(strRoleID)

	appCode := r.Header.Get(ctxs.UserAppCodeKey)

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
		w.Header().Set("Access-Control-Expose-Headers", ctxs.UserSetTokenKey)
		w.Header().Set(ctxs.UserSetTokenKey, resp.Token)
	}
	logx.WithContext(r.Context()).Infof("%s.CheckTokenWare ip:%v in.token=%s roleID：%v checkResp:%v",
		utils.FuncName(), strIP, strToken, strRoleID, utils.Fmt(resp))
	if roleID != 0 { //如果传了角色
		if !utils.SliceIn(roleID, resp.RoleIDs...) {
			err := errors.Parameter.AddMsgf("所选角色无权限")
			return nil, err
		}
	} else {
		roleID = resp.RoleIDs[0]
	}
	return &ctxs.UserCtx{
		IsOpen:     false,
		TenantCode: resp.TenantCode,
		AppCode:    appCode,
		UserID:     resp.UserID,
		RoleID:     roleID,
		IsAdmin:    resp.IsAdmin == def.True,
		IsAllData:  resp.IsAllData == def.True,
		IP:         strIP,
		Os:         r.Header.Get("User-Agent"),
	}, nil
}
