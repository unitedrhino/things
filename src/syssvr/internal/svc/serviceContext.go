package svc

import (
	"github.com/i-Things/things/shared/third/weixin"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/config"
	"github.com/i-Things/things/src/syssvr/internal/repo/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config        config.Config
	UserInfoModel mysql.UserInfoModel
	RoleInfoModle mysql.RoleInfoModel
	MenuInfoModle mysql.MenuInfoModel
	RoleMenuModle mysql.RoleMenuModel
	UserModel     mysql.UserModel
	RoleModel     mysql.RoleModel
	MenuModel     mysql.MenuModel
	WxMiniProgram *weixin.MiniProgram
	UserID        *utils.SnowFlake
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	ui := mysql.NewUserInfoModel(conn)
	ro := mysql.NewRoleInfoModel(conn)
	me := mysql.NewMenuInfoModel(conn)
	rom := mysql.NewRoleModel(conn, c.CacheRedis)
	mem := mysql.NewMenuModel(conn, c.CacheRedis)
	um := mysql.NewUserModel(conn, c.CacheRedis)

	WxMiniProgram := weixin.NewWexinMiniProgram(c.WexinMiniprogram, c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)

	return &ServiceContext{
		Config:        c,
		UserInfoModel: ui,
		UserModel:     um,
		RoleModel:     rom,
		MenuModel:     mem,
		RoleInfoModle: ro,
		MenuInfoModle: me,
		WxMiniProgram: WxMiniProgram,
		UserID:        UserID,
	}
}
