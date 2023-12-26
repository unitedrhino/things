package middleware

import (
	role "github.com/i-Things/things/src/syssvr/client/rolemanage"
	"net/http"
)

type CheckApiWareMiddleware struct {
	AuthRpc role.RoleManage
}

func NewCheckApiWareMiddleware() *CheckApiWareMiddleware {
	return &CheckApiWareMiddleware{}
}

func (m *CheckApiWareMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//userCtx := ctxs.GetUserCtx(r.Context())
		////校验 Casbin Rule
		//_, err := m.AuthRpc.RoleApiAuth(r.Context(), &user.RoleApiAuthReq{
		//	RoleID: userCtx.RoleID,
		//	Path:   r.URL.Path,
		//	Method: r.Method,
		//})
		//if err != nil {
		//	logx.WithContext(r.Context()).Errorf("%s.AuthApiCheck error=%s", utils.FuncName(), err)
		//	http.Error(w, "接口权限不足："+err.Error(), http.StatusUnauthorized)
		//	return
		//}
		next(w, r)
	}
}
