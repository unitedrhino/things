package clients

import (
	"database/sql"
	"github.com/i-Things/things/shared/conf"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
	//tdengine 的cgo模式，这个模式是最快的，需要可以打开
	//_ "github.com/taosdata/driver-go/v3/taosSql"
	_ "github.com/taosdata/driver-go/v3/taosWS"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type Td struct {
	*sql.DB
}

var (
	td   = Td{}
	once = sync.Once{}
)

func NewTDengine(DataSource conf.TSDB) (TD *Td, err error) {
	once.Do(func() {
		td.DB, err = sql.Open(DataSource.Driver, DataSource.DSN)
		if err != nil {
			return
		}
		td.DB.SetMaxIdleConns(200)
		_, err = td.Exec("create database if not exists ithings;")
	})
	if err != nil {
		logx.Errorf("tdengine 初始化失败,err:%v", err)
	}
	return &td, err
}
