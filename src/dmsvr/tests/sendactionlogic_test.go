package tests

import (
	"context"
	"fmt"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/dmclient"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"testing"
	"time"
)

func TestSendAction(t *testing.T) {
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}, Timeout: 1000 * 1000}))
	fmt.Println("TestSendAction")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	DeviceName := "test8"
	ProductID := "21CYs1k9YpG"
	_, err := client.SendAction(ctx, &dm.SendActionReq{
		ProductID:   ProductID,
		DeviceName:  DeviceName,
		ActionId:    "whistle",
		InputParams: "{\"time\":60,\"switch\":true}",
	})
	if err != nil {
		t.Errorf("%+v", errors.Fmt(err))
	}
}
