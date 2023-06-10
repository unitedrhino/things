package userHeader

const (
	UserUidKey      string = "iThings-uid"
	UserTokenKey    string = "iThings-token"
	UserSetTokenKey string = "iThings-set-token"
	MetadataKey     string = "iThings-meta"
)

type MetaField string

//注意：val值 必须是 首字母大写，其他小写
const (
	MetaFieldProjectID MetaField = "Ithings-Projectid" //meta里的项目权限控制ID字段（企业版功能）
)
