package middleware

import (
	role "gitee.com/i-Things/core/service/syssvr/client/rolemanage"
	user "gitee.com/i-Things/core/service/syssvr/client/usermanage"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/config"
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
		next(w, r)

		logx.WithContext(r.Context()).Infof("%s.Lifecycle.After", utils.FuncName())
	}
}
