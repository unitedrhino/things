package tests

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/usersvr/user"
	"github.com/go-things/things/src/usersvr/userclient"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
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
