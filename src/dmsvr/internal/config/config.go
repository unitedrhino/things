package config

import (
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/zrpc"
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
	Mqtt struct {
		Brokers []string //mqtt服务器节点
		User    string   //用户名
		Pass    string   `json:",optional"` //密码
	}
	Mongo struct {
		Url      string //mongodb连接串
		Database string //选择的数据库
	}
}
