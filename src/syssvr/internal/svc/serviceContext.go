package svc

import (
	"context"
	"github.com/casbin/casbin/v2"
	cas "github.com/i-Things/things/shared/casbin"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/config"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"os"
)

type ServiceContext struct {
	Config        config.Config
	WxMiniProgram *clients.MiniProgram
	UserID        *utils.SnowFlake
	Casbin        *casbin.Enforcer
	Store         kv.Store
}

func NewServiceContext(c config.Config) *ServiceContext {
	//conn := sqlx.NewMysql(c.Database.DSN)
	stores.InitConn(c.Database)

	// 自动迁移数据库
	db := stores.GetCommonConn(context.Background())
	errdb := db.AutoMigrate(&relationDB.SysUserInfo{})
	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}
	errdb = db.AutoMigrate(&relationDB.SysRoleInfo{})
	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}
	errdb = db.AutoMigrate(&relationDB.SysRoleMenu{})
	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}
	errdb = db.AutoMigrate(&relationDB.SysMenuInfo{})
	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}
	errdb = db.AutoMigrate(&relationDB.SysLoginLog{})
	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}
	errdb = db.AutoMigrate(&relationDB.SysOperLog{})
	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}
	errdb = db.AutoMigrate(&relationDB.SysApiInfo{})
	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}
	errdb = db.AutoMigrate(&relationDB.SysApiAuth{})

	if errdb != nil {
		logx.Error("failed to migrate database: %v", errdb)
		os.Exit(-1)
	}

	WxMiniProgram := clients.NewWxMiniProgram(context.Background(), c.WxMiniProgram, c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)
	dbRaw, err := db.DB()
	if err != nil {
		logx.Error("sys failed to  database err: %v", err)
	}
	ca := cas.NewCasbinWithRedisWatcher(dbRaw, c.Database.DBType, c.CacheRedis[0].RedisConf)
	store := kv.NewStore(c.CacheRedis)

	return &ServiceContext{
		Config:        c,
		WxMiniProgram: WxMiniProgram,
		UserID:        UserID,
		Casbin:        ca,
		Store:         store,
	}
}
