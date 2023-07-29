package svc

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/disvr/client/deviceinteract"
	"github.com/i-Things/things/src/disvr/didirect"
	"github.com/i-Things/things/src/stocksvr/internal/config"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_location"
	"github.com/i-Things/things/src/stocksvr/internal/models/stock_move"
	_ "github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Model     Model
	RpcClient RpcClient
}
type RpcClient struct {
	DeviceInteract deviceinteract.DeviceInteract
}
type Model struct {
	StockLocationModel stock_location.CacheModel
	StockMoveModel     stock_move.CacheModel
}

const postgresDriverName = "postgres"

func NewServiceContext(c config.Config) *ServiceContext {
	var deviceInteract deviceinteract.DeviceInteract
	if c.DiRpc.Enable {
		if c.DiRpc.Mode == conf.ClientModeGrpc {
			deviceInteract = deviceinteract.NewDeviceInteract(zrpc.MustNewClient(c.DiRpc.Conf))

		} else {
			deviceInteract = didirect.NewDeviceInteract(c.DiRpc.RunProxy)
		}
	}
	conn := sqlx.NewSqlConn(postgresDriverName, c.SqlConf.DataSource)
	return &ServiceContext{
		Config: c,
		Model: Model{
			StockLocationModel: stock_location.NewCacheModel(conn, c.CacheConf),
			StockMoveModel:     stock_move.NewCacheModel(conn, c.CacheConf),
		},
		RpcClient: RpcClient{DeviceInteract: deviceInteract},
	}
}
