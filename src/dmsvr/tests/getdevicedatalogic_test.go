package tests

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"testing"
)

var (
	productID  = "22ARVFc8Q0g"
	deviceName = "erw23"
)

func TestGetDeviceData(t *testing.T) {
	fmt.Println("GetDeviceData")
	client := dm.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}}))
	ctx := context.Background()
	{
		req := &dm.GetDeviceDataReq{
			Method:     "property",
			DeviceName: deviceName,
			ProductID:  productID,
			DataID:     "GPS_Info",
			Limit:      1,
			//TimeStart: time.Unix(1625013546,0).Unix(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceData(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dm.GetDeviceDataReq{
			Method:     "event",
			DeviceName: deviceName,
			ProductID:  productID,
			DataID:     "fesf",
			Limit:      10,
			//TimeStart: time.Unix(1625013546,0).Unix(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceData(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}

func TestGetDeviceDatas(t *testing.T) {
	fmt.Println("TestGetDeviceDatas")
	client := dm.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}}))
	ctx := context.Background()

	{
		req := &dm.GetDeviceDataReq{
			Method:     "property",
			DeviceName: deviceName,
			ProductID:  productID,
			DataID:     "",
			Limit:      1,
			//TimeStart: time.Unix(1625013546,0).Unix(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceData(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dm.GetDeviceDataReq{
			Method:     "event",
			DeviceName: deviceName,
			ProductID:  productID,
			DataID:     "fesf",
			Limit:      10,
			//TimeStart: time.Unix(1625013546,0).Unix(),
			//TimeEnd: time.Unix(1625223546,0).Unix(),
		}
		info, err := client.GetDeviceData(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}
