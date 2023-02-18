package svc

import (
	"github.com/i-Things/things/src/rulesvr/internal/config"
	"github.com/i-Things/things/src/rulesvr/internal/domain/scene"
	"github.com/i-Things/things/src/rulesvr/internal/repo/mysql"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Repo
}
type Repo struct {
	SceneRepo scene.Repo
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	SceneRepo := mysql.NewRuleSceneInfoModel(conn)
	return &ServiceContext{
		Config: c,
		Repo: Repo{
			SceneRepo: SceneRepo,
		},
	}
}
