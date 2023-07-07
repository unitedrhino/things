package svc

import (
	"context"
	"github.com/casbin/casbin/v2"
	cas "github.com/i-Things/things/shared/casbin"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/config"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	WxMiniProgram *clients.MiniProgram
	UserID        *utils.SnowFlake
	Casbin        *casbin.Enforcer
	Store         kv.Store
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Database.DSN)
	WxMiniProgram := clients.NewWxMiniProgram(context.Background(), c.WxMiniProgram, c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)
	db, _ := conn.RawDB()
	ca := cas.NewCasbinWithRedisWatcher(db, c.Database.DBType, c.CacheRedis[0].RedisConf)
	store := kv.NewStore(c.CacheRedis)

	return &ServiceContext{
		Config:        c,
		WxMiniProgram: WxMiniProgram,
		UserID:        UserID,
		Casbin:        ca,
		Store:         store,
	}
}
