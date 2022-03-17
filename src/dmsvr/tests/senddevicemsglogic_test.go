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

func TestSendDeviceMsg(t *testing.T) {
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}, Timeout: 1000 * 1000}))
	fmt.Println("TestSendDeviceMsg")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := client.SendDeviceMsg(ctx, &dm.SendDeviceMsgReq{
		ClientID:  "22ARVFc8Q0gerw23",
		Username:  "22ARVFc8Q0gerw23;12010126;5M5NH;1646579042",
		Topic:     "",
		Payload:   "",
		Timestamp: time.Now().Unix(),
		Action:    "connect",
	})
	if err != nil {
		t.Errorf("%+v", errors.Fmt(err))
	}
}
