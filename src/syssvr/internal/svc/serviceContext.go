package svc

import (
	"context"
	"github.com/casbin/casbin/v2"
	cas "github.com/i-Things/things/shared/casbin"
	"github.com/i-Things/things/shared/clients"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/config"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	UserInfoModel mysql.SysUserInfoModel
	RoleInfoModel mysql.SysRoleInfoModel
	MenuInfoModel mysql.SysMenuInfoModel
	RoleMenuModel mysql.SysRoleMenuModel
	UserModel     mysql.UserModel
	RoleModel     mysql.RoleModel
	MenuModel     mysql.MenuModel
	WxMiniProgram *clients.MiniProgram
	UserID        *utils.SnowFlake
	LogLoginModel mysql.SysLoginLogModel
	LogOperModel  mysql.SysOperLogModel
	LogModel      mysql.LogModel
	ApiModel      mysql.SysApiInfoModel
	Casbin        *casbin.Enforcer
	Store         kv.Store
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Database.DSN)
	ui := mysql.NewSysUserInfoModel(conn)
	ro := mysql.NewSysRoleInfoModel(conn)
	me := mysql.NewSysMenuInfoModel(conn)
	ll := mysql.NewSysLoginLogModel(conn)
	lo := mysql.NewSysOperLogModel(conn)
	l := mysql.NewLogModel(conn)
	api := mysql.NewSysApiInfoModel(conn)
	rom := mysql.NewRoleModel(conn, c.CacheRedis)
	mem := mysql.NewMenuModel(conn, c.CacheRedis)
	um := mysql.NewUserModel(conn, c.CacheRedis)

	WxMiniProgram := clients.NewWxMiniProgram(context.Background(), c.WxMiniProgram, c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)
	db, _ := conn.RawDB()
	ca := cas.NewCasbinWithRedisWatcher(db, c.Database.DBType, c.CacheRedis[0].RedisConf)
	store := kv.NewStore(c.CacheRedis)

	return &ServiceContext{
		Config:        c,
		UserInfoModel: ui,
		UserModel:     um,
		RoleModel:     rom,
		MenuModel:     mem,
		RoleInfoModel: ro,
		MenuInfoModel: me,
		WxMiniProgram: WxMiniProgram,
		UserID:        UserID,
		LogLoginModel: ll,
		LogOperModel:  lo,
		LogModel:      l,
		ApiModel:      api,
		Casbin:        ca,
		Store:         store,
	}
}
