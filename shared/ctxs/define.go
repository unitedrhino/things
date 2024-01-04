package ctxs

import "strings"

const (
	UserInfoKey       string = "ithings-user"
	UserTokenKey      string = "ithings-token"
	UserAppCodeKey    string = "ithings-app-code"    //用户正在访问的app
	UserTenantCodeKey string = "ithings-tenant-code" //用户租户号

	UserRoleKey     string = "ithings-user-role"
	UserProjectID   string = "ithings-project-id"
	UserSetTokenKey string = "ithings-set-token"
	MetadataKey     string = "ithings-meta"
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
