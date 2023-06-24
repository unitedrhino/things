package conf

const (
	Mysql  = "mysql"
	Pgsql  = "pgsql"
	Sqlite = "sqlite"
)

type Database struct {
	DBType string `json:",default=mysql,options=mysql|pgsql"` //
	DSN    string `json:""`                                   //dsn
}
