package svc

import (
	"gitee.com/godLei6/things/shared/utils"
	"gitee.com/godLei6/things/src/dcsvr/internal/config"
	"gitee.com/godLei6/things/src/dcsvr/model"
	"gitee.com/godLei6/things/src/dmsvr/dmclient"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	GroupInfo model.GroupInfoModel
	GroupMember model.GroupMemberModel
	DcDB   model.DmModel
	GroupID   *utils.SnowFlake
	Dmsvr   dmclient.Dm
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	gi := model.NewGroupInfoModel(conn, c.CacheRedis)
	gm := model.NewGroupMemberModel(conn, c.CacheRedis)
	dc := model.NewDcModel(conn,c.CacheRedis)
	GroupID := utils.NewSnowFlake(c.NodeID)
	dm := dmclient.NewDm(zrpc.MustNewClient(c.DmRpc.Conf))

	return &ServiceContext{
		Config: c,
		GroupInfo: gi,
		GroupMember: gm,
		DcDB: dc,
		GroupID: GroupID,
		Dmsvr:dm,
	}
}
