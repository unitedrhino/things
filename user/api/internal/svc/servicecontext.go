package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"yl/user/model"
	"yl/user/api/internal/config"
	"yl/user/api/internal/middleware"
)

type ServiceContext struct {
	Config    			config.Config
	Usercheck 			rest.Middleware
	UserInfoModel 		model.UserInfoModel
	UserCoreModel      	model.UserCoreModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	ui := model.NewUserInfoModel(conn)
	uc := model.NewUserCoreModel(conn)
	return &ServiceContext{
		Config:    c,
		Usercheck: middleware.NewUsercheckMiddleware().Handle,
		UserInfoModel: ui,
		UserCoreModel:uc,
	}
}
