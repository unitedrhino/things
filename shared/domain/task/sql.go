package task

const (
	SqlTypeNormal = "normal"
	SqlTypeJs     = "js"
)
const (
	SqlEnvDsn    = "dsn"
	SqlEnvDBType = "dbType"
)

type Sql struct {
	Type   string `json:"type"` // sql执行的类型  normal(支持环境变量的直接执行方式),js(js脚本执行)
	Sql    string `json:"sql"`  // 执行的sql
	Script string `json:"script"`
	//需要使用的环境变量 dsn:如果填写,默认使用该dsn来连接,dbType:mysql|pgsql 默认是mysql
	Env map[string]any `json:"env"`
}
