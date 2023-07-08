package stores

import (
	"context"
	"github.com/glebarez/sqlite"
	"github.com/i-Things/things/shared/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	commonConn *gorm.DB
	once       sync.Once
	tenantConn sync.Map
)

func InitConn(database conf.Database) {
	var err error
	once.Do(func() {
		switch database.DBType {
		case conf.Mysql:
			commonConn, err = gorm.Open(mysql.Open(database.DSN), &gorm.Config{})
		case conf.Pgsql:
			commonConn, err = gorm.Open(postgres.Open(database.DSN), &gorm.Config{})
		case conf.Sqlite:
			commonConn, err = gorm.Open(sqlite.Open(database.DSN), &gorm.Config{})
		}
		logx.Must(err)
	})
	return
}

// 获取租户连接  传入context或db连接 如果传入的是db连接则直接返回db
func GetTenantConn(in any) *gorm.DB {
	if db, ok := in.(*gorm.DB); ok {
		return db
	}
	return commonConn.WithContext(in.(context.Context)).Debug()
}

// 获取公共连接 传入context或db连接 如果传入的是db连接则直接返回db
func GetCommonConn(in any) *gorm.DB {
	if db, ok := in.(*gorm.DB); ok {
		return db
	}
	return commonConn.WithContext(in.(context.Context)).Debug()
}
