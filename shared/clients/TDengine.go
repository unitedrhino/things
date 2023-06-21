package clients

import (
	"database/sql"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
	"sync"
)

type Td struct {
	Dsn string
	*sql.DB
}

var (
	td   = Td{}
	once = sync.Once{}
)

func NewTDengine(DataSource string) (TD *Td, err error) {
	once.Do(func() {
		td.DB, err = sql.Open("taosRestful", DataSource)
		if err != nil {
			return
		}
		td.Dsn = DataSource
		_, err = td.Exec("create database if not exists iThings;")
	})
	return &td, err
}
