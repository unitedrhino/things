package tests

import (
	"context"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/usersvr/user"
	"gitee.com/godLei6/things/src/usersvr/userclient"
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc"
	"testing"
)

func TestGetUserCoreList(t *testing.T) {
	t.Log("TestGetUserCore")
	client := userclient.NewUser(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "user.rpc",
	}}))
	ctx := context.Background()

	req := &user.GetUserCoreListReq{
		Page: &user.PageInfo{
			Page:       1,
			PageSize:   20,
			SearchKey:  "fwef",
			SearchType: "agrhgsrgr",
		},
	}
	info, err := client.GetUserCoreList(ctx, req)
	t.Log(req, info)
	if err != nil {
		t.Errorf("%+v", errors.Fmt(err))
	}

}
