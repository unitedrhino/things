package svc

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dcsvr/internal/config"
	"github.com/i-Things/things/src/dcsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/dmdirect"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	GroupInfo   mysql.GroupInfoModel
	GroupMember mysql.GroupMemberModel
	DcDB        mysql.DmModel
	GroupID     *utils.SnowFlake
	Dmsvr       dm.Dm
}

func NewServiceContext(c config.Config) *ServiceContext {
	var dr dm.Dm
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	gi := mysql.NewGroupInfoModel(conn, c.CacheRedis)
	gm := mysql.NewGroupMemberModel(conn, c.CacheRedis)
	dc := mysql.NewDcModel(conn, c.CacheRedis)
	nodeId := utils.GetNodeID(c.CacheRedis, c.Name)
	GroupID := utils.NewSnowFlake(nodeId)
	if c.DmRpc.Mode == conf.ClientModeGrpc {
		dr = dm.NewDm(zrpc.MustNewClient(c.DmRpc.Conf))
	} else {
		dr = dmdirect.NewDm(nil)
	}
	return &ServiceContext{
		Config:      c,
		GroupInfo:   gi,
		GroupMember: gm,
		DcDB:        dc,
		GroupID:     GroupID,
		Dmsvr:       dr,
	}
}
