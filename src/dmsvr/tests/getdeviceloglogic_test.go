package tests

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

func TestGetDeviceLog(t *testing.T) {
	fmt.Println("GetDeviceLog")
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}}))
	ctx := context.Background()
	productID := "21CYs1k9YpG"
	deviceName := "test8"
	{
		req := &dm.GetDeviceLogReq{
			Method:     "property",
			DeviceName: deviceName,
			ProductID:  productID,
			DataID:     "GPS_Info",
			Limit:      1,
			//TimeStart: time.Unix(1625013546,0).Unix(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceLog(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dm.GetDeviceLogReq{
			Method:     "event",
			DeviceName: deviceName,
			ProductID:  productID,
			DataID:     "fesf",
			Limit:      10,
			//TimeStart: time.Unix(1625013546,0).Unix(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceLog(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}
