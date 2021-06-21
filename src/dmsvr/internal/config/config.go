package config

import (
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	NodeID int64		//节点id
	Mysql struct {
		DataSource string
	}
	CacheRedis cache.ClusterConf
	Kafka      struct{
		Brokers		[]string	//kafka的节点
		Group 		string		//kafka的分组
	}
}

