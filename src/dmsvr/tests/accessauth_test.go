package tests

import (
	"context"
	"fmt"
	"github.com/go-things/things/shared/def"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/dmclient"
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc"
	"testing"
)

func TestAccessAuth(t *testing.T) {
	fmt.Println("TestLoginAuth")
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}}))
	ctx := context.Background()
	productID := "21CYs1k9YpG"
	deviceName := "test8"
	clientID := fmt.Sprintf("%s%s", productID, deviceName)
	//生成 MQTT 的 username 部分, 格式为 ${clientid};${sdkappid};${connid};${expiry}
	userName := fmt.Sprintf("%s;12010126;fawef;1822730956", clientID)

	topics := []string{
		"$thing/up/property/%s/%s",
		//"$thing/down/property/%s/%s",
		//"ota/report/%s/%s",
		//"broadcast/rxd/%s/%s",
		"%s/%s/control",
	}
	for _, v := range topics {
		req := &dm.AccessAuthReq{
			Username: userName, //用户名
			Topic:    fmt.Sprintf(v, productID, deviceName),
			Access:   def.PUB,
			ClientID: clientID,
			Ip:       "192.168.1.2", //访问的ip地址
		}
		info, err := client.AccessAuth(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}
