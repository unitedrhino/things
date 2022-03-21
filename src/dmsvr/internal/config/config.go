package config

import (
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/src/dmsvr/internal/repo/third"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	NodeID int64 //节点id
	Mysql  struct {
		DataSource string
	}
	CacheRedis cache.ClusterConf
	Kafka      struct {
		Brokers []string //kafka的节点
		Group   string   //kafka的分组
	}
	DevClient third.DevClientConf
	Mongo     struct {
		Url      string //mongodb连接串
		Database string //选择的数据库
	}
	InnerLink InnerLinkConf //和things内部交互的设置
}

type InnerLinkConf struct {
	Nats conf.NatsConf
}
