package conf

const (
	Mysql  = "mysql"
	Pgsql  = "pgsql"
	Sqlite = "sqlite"
)

type Database struct {
	DBType      string `json:",default=mysql,options=mysql|pgsql"` //
	IsInitTable bool   `json:",default=true"`
	DSN         string `json:""` //dsn
}

// 时序数据库（Time Series Database）
type TSDB struct {
	DBType string `json:",default=tdengine,options=tdengine"`         //
	Driver string `json:",default=taosWS,options=taosRestful|taosWS"` //
	DSN    string `json:""`                                           //dsn
}
