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
	WxMiniProgram *weixin.MiniProgram
	UserID        *utils.SnowFlake
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	ui := mysql.NewUserInfoModel(conn)
	ro := mysql.NewRoleInfoModel(conn)
	um := mysql.NewUserModel(conn, c.CacheRedis)
	WxMiniProgram := weixin.NewWexinMiniProgram(c.WexinMiniprogram, c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)

	return &ServiceContext{
		Config:        c,
		UserInfoModel: ui,
		UserModel:     um,
		RoleInfoModle: ro,
		WxMiniProgram: WxMiniProgram,
		UserID:        UserID,
	}
}
