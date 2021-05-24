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
func TestManageProduct(t *testing.T) {
	client := dmclient.NewDm(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key: "dm.rpc",
	}}))
	 //l := CheckTokenLogic{}
	 //resp,err := l.CheckToken(&user.CheckTokenReq{
	 //	Token: "123123",
	 //})
	 //t.Errorf("TestCheckToken|resp=%#v|err=%#v\n",resp,err)
	fmt.Println("TestCheckToken")
	ctx := context.Background()
	ProductName := "test5"
	info,err := client.ManageProduct(ctx,&dm.ManageProductReq{
		Opt: dm.OptType_ADD,
		Info: &dm.ProductInfo{
			ProductName: ProductName,
		},
	})
	if err != nil {
		t.Errorf("%+v",errors.Fmt(err))
	}
	fmt.Println(info)
	if info.ProductName != ProductName{
		t.Errorf("ProductName not succ:%s",info.ProductName)
	}
	_,err = client.ManageProduct(ctx,&dm.ManageProductReq{
		Opt: dm.OptType_ADD,
		Info: &dm.ProductInfo{
			ProductName: ProductName,
		},
	})
	if !errors.Cmp(err,errors.Duplicate){
		t.Errorf("need duplicate err")
	}
	info,err = client.ManageProduct(ctx,&dm.ManageProductReq{
		Opt: dm.OptType_MODIFY,
		Info: &dm.ProductInfo{
			ProductID: info.ProductID,
			ProductName: ProductName+"1",
		},
	})
	if err != nil {
		t.Errorf("%+v",errors.Fmt(err))
	}
	if info.ProductName != (ProductName+"1"){
		t.Errorf("%+v", info)
	}

	info,err = client.ManageProduct(ctx,&dm.ManageProductReq{
		Opt: dm.OptType_DEL,
		Info: &dm.ProductInfo{
			ProductID: info.ProductID,
		},
	})
	if err != nil {
		t.Errorf("%+v",errors.Fmt(err))
	}
}
