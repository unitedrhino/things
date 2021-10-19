package logic

import (
	"context"
	"fmt"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dcsvr/dc"
	"gitee.com/godLei6/things/src/dcsvr/dcclient"
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc"
	"testing"
	"time"
)

func TestSendAction(t *testing.T) {
	client := dcclient.NewDc(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dc.rpc",
	},Timeout: 1000*1000}))
	fmt.Println("TestSendAction")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	DeviceName := "test8"
	ProductID := "21CYs1k9YpG"
	_, err := client.SendAction(ctx, &dc.SendActionReq{
		MemberID: "1699809227385606144",
		MemberType: 2,
		ProductID            :ProductID,
		DeviceName           :DeviceName,
		ActionId             :"whistle",
		InputParams          :"{\"time\":60,\"switch\":true}",
	})
	if err != nil {
		t.Errorf("%+v", errors.Fmt(err))
	}
}
