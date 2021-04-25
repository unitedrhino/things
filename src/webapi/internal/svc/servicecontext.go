package svc

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/zrpc"
	"yl/src/user/model"
	"yl/src/user/userclient"
	"yl/src/webapi/internal/config"
	"yl/src/webapi/internal/middleware"
)

type ServiceContext struct {
	Config        config.Config
	CheckToken    rest.Middleware
	Record		  rest.Middleware
	UserInfoModel model.UserInfoModel
	UserCoreModel model.UserCoreModel
	UserRpc       userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	ui := model.NewUserInfoModel(conn,c.CacheRedis)
	uc := model.NewUserCoreModel(conn,c.CacheRedis)
	ur := userclient.NewUser(zrpc.MustNewClient(c.UserRpc))
	return &ServiceContext{
		Config:        c,
		CheckToken:    middleware.NewCheckTokenMiddleware(ur).Handle,
		Record:		   middleware.NewRecordMiddleware().Handle,
		UserInfoModel: ui,
		UserCoreModel: uc,
		UserRpc:       ur,
	}
}
