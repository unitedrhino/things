package svc

import (
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dcsvr/internal/config"
	"gitee.com/godLei6/things/src/dcsvr/model"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	GroupInfo model.GroupInfoModel
	GroupMember model.GroupMemberModel
	DcDB   model.DmModel
	GroupID   *utils.SnowFlake
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	gi := model.NewGroupInfoModel(conn, c.CacheRedis)
	gm := model.NewGroupMemberModel(conn, c.CacheRedis)
	dc := model.NewDcModel(conn,c.CacheRedis)
	GroupID := utils.NewSnowFlake(c.NodeID)
	return &ServiceContext{
		Config: c,
		GroupInfo: gi,
		GroupMember: gm,
		DcDB: dc,
		GroupID: GroupID,
	}
}
