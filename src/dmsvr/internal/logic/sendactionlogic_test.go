package logic

import (
	"context"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/dmclient"
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc"
	"testing"
)

func TestSendAction(t *testing.T) {
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}}))
	fmt.Println("TestSendAction")
	ctx := context.Background()
	DeviceName := "test8"
	ProductID := "21CYs1k9YpG"
	_, err := client.SendAction(ctx, &dm.SendActionReq{
		ProductID            :ProductID,
		DeviceName           :DeviceName,
		ActionId             :"action id",
		InputParams          :"{\"a\":\"hello world\"}",
	})
	if err != nil {
		t.Errorf("%+v", errors.Fmt(err))
	}
}
