// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conf

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	"ariga.io/entcache"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zeromicro/go-zero/core/logx"
	redis2 "github.com/zeromicro/go-zero/core/stores/redis"
)

// DatabaseConf stores database configurations.
type DatabaseConf struct {
	Host         string
	Port         int
	Username     string `json:",default=root"`
	Password     string `json:",optional"`
	DBName       string `json:",default=simple_admin"`
	SSLMode      string `json:",optional"`
	Type         string `json:",default=mysql,options=[mysql,postgres,sqlite3]"`
	MaxOpenConns int    `json:",optional,default=100"`
	Debug        bool   `json:",optional,default=false"`
	CacheTime    int    `json:",optional,default=10"`
	DBPath       string `json:",optional"`
	MysqlConfig  string `json:",optional"`
	PGConfig     string `json:",optional"`
	SqliteConfig string `json:",optional"`
}

// NewCacheDriver returns an Ent driver with cache.
func (c DatabaseConf) NewCacheDriver(redisConf redis2.RedisConf) *entcache.Driver {
	db, err := sql.Open(c.Type, c.GetDSN())
	logx.Must(err)

	db.SetMaxOpenConns(c.MaxOpenConns)
	driver := entsql.OpenDB(c.Type, db)

	rdb := redis.NewClient(&redis.Options{Addr: redisConf.Host})

	cacheDrv := entcache.NewDriver(
		driver,
		entcache.TTL(time.Duration(c.CacheTime)*time.Second),
		entcache.Levels(
			entcache.NewLRU(256),
			entcache.NewRedis(rdb),
		),
	)

	return cacheDrv
}

// NewNoCacheDriver returns an Ent driver without cache.
func (c DatabaseConf) NewNoCacheDriver() *entsql.Driver {
	db, err := sql.Open(c.Type, c.GetDSN())
	logx.Must(err)

	db.SetMaxOpenConns(c.MaxOpenConns)
	driver := entsql.OpenDB(c.Type, db)

	return driver
}

// MysqlDSN returns mysql DSN.
func (c DatabaseConf) MysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True%s", c.Username, c.Password, c.Host, c.Port, c.DBName, c.MysqlConfig)
}

// PostgresDSN returns Postgres DSN.
func (c DatabaseConf) PostgresDSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s%s", c.Username, c.Password, c.Host, c.Port, c.DBName,
		c.SSLMode, c.PGConfig)
}

// SqliteDSN returns Sqlite DSN.
func (c DatabaseConf) SqliteDSN() string {
	if c.DBPath == "" {
		logx.Must(errors.New("the database file path cannot be empty"))
	}

	if _, err := os.Stat(c.DBPath); os.IsNotExist(err) {
		f, err := os.OpenFile(c.DBPath, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			logx.Must(fmt.Errorf("failed to create SQLite database file %q", c.DBPath))
		}
		if err := f.Close(); err != nil {
			logx.Must(fmt.Errorf("failed to create SQLite database file %q", c.DBPath))
		}
	} else {
		if err := os.Chmod(c.DBPath, 0660); err != nil {
			logx.Must(fmt.Errorf("unable to set permission code on %s: %v", c.DBPath, err))
		}
	}

	return fmt.Sprintf("file:%s?_busy_timeout=100000&_fk=1%s", c.DBPath, c.SqliteConfig)
}

// GetDSN returns DSN according to the database type.
func (c DatabaseConf) GetDSN() string {
	switch c.Type {
	case "mysql":
		return c.MysqlDSN()
	case "postgres":
		return c.PostgresDSN()
	case "sqlite3":
		return c.SqliteDSN()
	default:
		return "mysql"
	}
}
