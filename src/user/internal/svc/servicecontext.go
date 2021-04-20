package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"yl/shared/third/weixin"
	"yl/shared/utils"

	"yl/src/user/internal/config"
	//"yl/shared/third/weixin"
	"yl/src/user/model"
)

type ServiceContext struct {
	Config        config.Config
	UserInfoModel model.UserInfoModel
	UserCoreModel model.UserCoreModel
	UserModel     model.UserModel
	WxMiniProgram *weixin.MiniProgram
	UserID   	  *utils.SnowFlake
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	ui := model.NewUserInfoModel(conn,c.CacheRedis)
	uc := model.NewUserCoreModel(conn,c.CacheRedis)
	um := model.NewUserModel(conn,c.CacheRedis)
	WxMiniProgram :=  weixin.NewWexinMiniProgram(c.WexinMiniprogram,c.CacheRedis)
	UserID := utils.NewSnowFlake(c.NodeID)


	return &ServiceContext{
		Config: c,
		UserInfoModel: ui,
		UserCoreModel: uc,
		UserModel: um,
		WxMiniProgram:WxMiniProgram,
		UserID:UserID,
	}
}
