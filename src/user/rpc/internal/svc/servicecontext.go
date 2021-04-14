package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"yl/shared/third/weixin"

	//"yl/shared/third/weixin"
	"yl/src/user/model"
	"yl/src/user/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	UserInfoModel model.UserInfoModel
	UserCoreModel model.UserCoreModel
	WxMiniProgram *weixin.MiniProgram
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	ui := model.NewUserInfoModel(conn,c.CacheRedis)
	uc := model.NewUserCoreModel(conn,c.CacheRedis)
	WxMiniProgram :=  weixin.NewWexinMiniProgram(c.WexinMiniprogram,c.CacheRedis)
	return &ServiceContext{
		Config: c,
		UserInfoModel: ui,
		UserCoreModel: uc,
		WxMiniProgram:WxMiniProgram,
	}
}
