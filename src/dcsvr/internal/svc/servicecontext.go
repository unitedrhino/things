package svc

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dcsvr/internal/config"
	"github.com/i-Things/things/src/dcsvr/model"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	GroupInfo   model.GroupInfoModel
	GroupMember model.GroupMemberModel
	DcDB        model.DmModel
	GroupID     *utils.SnowFlake
	Dmsvr       dm.Dm
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	gi := model.NewGroupInfoModel(conn, c.CacheRedis)
	gm := model.NewGroupMemberModel(conn, c.CacheRedis)
	dc := model.NewDcModel(conn, c.CacheRedis)
	GroupID := utils.NewSnowFlake(c.NodeID)
	dm := dm.NewDm(zrpc.MustNewClient(c.DmRpc.Conf))

	return &ServiceContext{
		Config:      c,
		GroupInfo:   gi,
		GroupMember: gm,
		DcDB:        dc,
		GroupID:     GroupID,
		Dmsvr:       dm,
	}
}
