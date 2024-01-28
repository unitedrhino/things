package cache

import (
	"github.com/zeromicro/go-zero/core/stores/kv"
)

type ApiAuth struct {
	store kv.Store
}

func NewApiAuth(store kv.Store) *ApiAuth {
	return &ApiAuth{
		store: store,
	}
}

//func (c *ApiAuth) GenKey(role, codeID string) string {
//	return "captcha:" + Type + ":" + codeID
//}
//
//func (c *ApiAuth) Reload(ctx context.Context, roleID int64, tenantCode, appCode, moduleCode string) error {
//	ctxs.GetUserCtx(ctx).AllTenant = true
//	defer func() {
//		ctxs.GetUserCtx(ctx).AllTenant = false
//	}()
//	list, err := relationDB.NewRoleApiRepo(ctx).FindByFilter(ctx, relationDB.RoleApiFilter{
//		TenantCode: tenantCode,
//		AppCode:    appCode,
//		AccessCode: moduleCode,
//		RoleIDs:    []int64{roleID},
//	}, nil)
//	if err != nil {
//		return err
//	}
//
//}

//func (c *ApiAuth) Verify(ctx context.Context, Type, codeID, code string) string {
//	key := c.GenKey(Type, codeID)
//	val, err := c.store.GetCtx(ctx, key)
//	if err != nil || val == "" {
//		return ""
//	}
//	//如果验证码存在，则删除验证码
//	c.store.DelCtx(ctx, key)
//	body := map[string]string{}
//	json.Unmarshal([]byte(val), &body)
//	if body["code"] == code {
//		if body["account"] == "" {
//			return " "
//		}
//		return body["account"]
//	}
//	return ""
//}
//
//func (c *ApiAuth) Store(ctx context.Context, Type, codeID, code string, account string, expire int64) error {
//	body := map[string]interface{}{
//		"code":    code,
//		"account": account,
//	}
//	bodytStr, _ := json.Marshal(body)
//	return c.store.SetexCtx(ctx, c.GenKey(Type, codeID), string(bodytStr), int(expire))
//}
