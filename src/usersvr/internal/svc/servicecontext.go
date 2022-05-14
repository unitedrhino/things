package svc

import (
	"github.com/i-Things/things/shared/third/weixin"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/usersvr/internal/repo/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"github.com/i-Things/things/src/usersvr/internal/config"
)

type ServiceContext struct {
	Config        config.Config
	UserInfoModel mysql.UserInfoModel
	UserCoreModel mysql.UserCoreModel
	UserModel     mysql.UserModel
	WxMiniProgram *weixin.MiniProgram
	UserID        *utils.SnowFlake
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	ui := mysql.NewUserInfoModel(conn, c.CacheRedis)
	uc := mysql.NewUserCoreModel(conn, c.CacheRedis)
	um := mysql.NewUserModel(conn, c.CacheRedis)
	WxMiniProgram := weixin.NewWexinMiniProgram(c.WexinMiniprogram, c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	UserID := utils.NewSnowFlake(nodeId)

	return &ServiceContext{
		Config:        c,
		UserInfoModel: ui,
		UserCoreModel: uc,
		UserModel:     um,
		WxMiniProgram: WxMiniProgram,
		UserID:        UserID,
	}
}
