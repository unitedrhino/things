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
	dbType     string //数据库类型
)

func InitConn(database conf.Database) {
	var err error
	once.Do(func() {
		commonConn, err = GetConn(database)
		logx.Must(err)
	})
	return
}
func GetConn(database conf.Database) (conn *gorm.DB, err error) {
	dbType = database.DBType
	switch database.DBType {
	case conf.Pgsql:
		conn, err = gorm.Open(postgres.Open(database.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	case conf.Sqlite:
		conn, err = gorm.Open(sqlite.Open(database.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	default:
		conn, err = gorm.Open(mysql.Open(database.DSN), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	}
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
	if in == nil {
		return commonConn.Debug()
	}
	return commonConn.WithContext(in.(context.Context)).Debug()
}
