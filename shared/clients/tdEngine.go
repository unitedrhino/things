package clients

import (
	"database/sql"
	"github.com/i-Things/things/shared/conf"
	_ "github.com/taosdata/driver-go/v3/taosRestful"
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
		driver := "taosWS"
		if DataSource.Driver == "taosRestful" {
			driver = "taosRestful"
		}
		td.DB, err = sql.Open(driver, DataSource.DSN)
		if err != nil {
			return
		}
		_, err = td.Exec("create database if not exists ithings;")
	})
	if err != nil {
		logx.Errorf("tdengine 初始化失败,err:%v", err)
	}
	return &td, err
}
