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

func TestGetDescribeDeviceLog(t *testing.T) {
	fmt.Println("TestGetDescribeDeviceLog")
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}}))
	ctx := context.Background()
	productID := "22ARVFc8Q0g"
	deviceName := "erw23"
	{
		req := &dm.GetDeviceDescribeLogReq{
			DeviceName: deviceName,
			ProductID:  productID,
			Limit:      100,
			//TimeStart: time.Unix(1625013546,0).Unix(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceDescribeLog(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dm.GetDeviceDescribeLogReq{
			DeviceName: deviceName,
			ProductID:  productID,
			Limit:      10,
			TimeStart:  time.Unix(1625013546, 0).UnixMilli(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceDescribeLog(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}
