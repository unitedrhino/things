package svc

import (
	"context"
	"github.com/i-Things/things/src/rulesvr/internal/config"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/repo/cache"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"os"
)

type ServiceContext struct {
	Config config.Config
	Repo
}
type Repo struct {
	SceneRepo       scene.Repo
	SceneDeviceRepo scene.DeviceRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	SceneRepo := mysql.NewRuleSceneInfoModel(conn)
	sceneDevice := cache.NewSceneDeviceRepo(SceneRepo)
	err := sceneDevice.Init(context.TODO())
	if err != nil {
		logx.Error("设备场景数据初始化失败 err:", err)
		os.Exit(-1)
	}
	return &ServiceContext{
		Config: c,
		Repo: Repo{
			SceneRepo:       SceneRepo,
			SceneDeviceRepo: sceneDevice,
		},
	}
}
