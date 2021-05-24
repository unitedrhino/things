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

func TestManageDevice(t *testing.T) {
	fmt.Println("TestManageDevice")
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key: "dm.rpc",
	}}))
	ctx := context.Background()
	Name := "test7"
	info,err := client.ManageDevice(ctx,&dm.ManageDeviceReq{
		Opt: dm.OptType_ADD,
		Info: &dm.DeviceInfo{
			ProductID: 1700642615085502464,
			DeviceName: Name,
		},
	})
	if err != nil {
		t.Errorf("%+v",errors.Fmt(err))
	}
	fmt.Println(info)
	if info.DeviceName != Name{
		t.Errorf("DeviceName not succ:%s",info.DeviceName)
	}
	_,err = client.ManageDevice(ctx,&dm.ManageDeviceReq{
		Opt: dm.OptType_ADD,
		Info: &dm.DeviceInfo{
			ProductID: 123,
			DeviceName: Name,
		},
	})
	if !errors.Cmp(err,errors.Duplicate){
		t.Errorf("need duplicate err")
	}
	info,err = client.ManageDevice(ctx,&dm.ManageDeviceReq{
		Opt: dm.OptType_MODIFY,
		Info: &dm.DeviceInfo{
			DeviceID: info.DeviceID,
			DeviceName: Name+"1",
		},
	})
	if err != nil {
		t.Errorf("%+v",errors.Fmt(err))
	}
	if info.DeviceName != (Name+"1"){
		t.Errorf("%+v", info)
	}

	info,err = client.ManageDevice(ctx,&dm.ManageDeviceReq{
		Opt: dm.OptType_DEL,
		Info: &dm.DeviceInfo{
			DeviceID: info.DeviceID,
		},
	})
	if err != nil {
		t.Errorf("%+v",errors.Fmt(err))
	}
}
