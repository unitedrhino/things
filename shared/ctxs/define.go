package ctxs

import "strings"

const (
	UserInfoKey       string = "iThings-user"
	UserTokenKey      string = "iThings-token"
	UserAppCodeKey    string = "iThings-app-code"    //用户正在访问的app
	UserTenantCodeKey string = "iThings-tenant-code" //用户租户号

	UserRoleKey     string = "iThings-user-role"
	UserSetTokenKey string = "iThings-set-token"
	MetadataKey     string = "iThings-meta"
)

type MetaField string

// 注意：val值 必须是 首字母大写，其他小写
const (
	MetaFieldProjectID MetaField = "Ithings-Project-Id" //meta里的项目权限控制ID字段（企业版功能）
)

var HttpAllowHeader string

func init() {
	HttpAllowHeader = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With," + strings.Join(ContextKeys, ",")
}
