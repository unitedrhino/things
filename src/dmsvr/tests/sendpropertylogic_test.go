package tests

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/dmclient"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"testing"
	"time"
)

func TestSendProperty(t *testing.T) {
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}, Timeout: 1000 * 1000}))
	fmt.Println("TestSendProperty")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	DeviceName := "test1"
	ProductID := "22BIUqIZSve"
	resp, err := client.SendProperty(ctx, &dm.SendPropertyReq{
		ProductID:  ProductID,
		DeviceName: DeviceName,
		Data:       "{\"time\":60,\"switch\":true}",
	})
	fmt.Println(resp, err)
	if err != nil {
		t.Errorf("%+v", errors.Fmt(err))
	}
}
