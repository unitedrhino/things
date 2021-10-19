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

func TestGetDeviceInfo(t *testing.T) {
	fmt.Println("GetDeviceInfo")
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dm.rpc",
	}}))
	ctx := context.Background()
	productID := "21CYs1k9YpG"
	deviceName := "test8"
	{
		req := &dm.GetDeviceInfoReq{
			DeviceName: deviceName, //设备名 为空时获取产品id下的所有设备信息
			ProductID:  productID,  //产品id
			//Page       : productID,//分页信息 只获取一个则不填
		}
		info, err := client.GetDeviceInfo(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dm.GetDeviceInfoReq{
			DeviceName: "",        //设备名 为空时获取产品id下的所有设备信息
			ProductID:  productID, //产品id
			Page: &dm.PageInfo{
				Page:     1,
				PageSize: 20,
			}, //分页信息 只获取一个则不填
		}
		info, err := client.GetDeviceInfo(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dm.GetDeviceInfoReq{
			DeviceName: "",    //设备名 为空时获取产品id下的所有设备信息
			ProductID:  "123", //产品id
			Page: &dm.PageInfo{
				Page:     1,
				PageSize: 20,
			}, //分页信息 只获取一个则不填
		}
		info, err := client.GetDeviceInfo(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}
