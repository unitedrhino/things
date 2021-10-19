package logic

import (
	"context"
	"fmt"
	"gitee.com/godLei6/things/shared/def"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dcsvr/dcclient"
	"github.com/tal-tech/go-zero/core/discov"
	"github.com/tal-tech/go-zero/zrpc"
	"testing"
)

func TestManageGroupInfo(t *testing.T) {
	fmt.Println("ManageGroupInfo")
	client := dcclient.NewDc(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dc.rpc",
	}}))
	ctx := context.Background()
	{
		req := &dcclient.ManageGroupInfoReq{
			Opt: def.OPT_ADD,
			Info: &dcclient.GroupInfo{
				Name        :  "测试组1",              //组名
				Uid        :    1699809227385606144,             //管理员用户id
			},
		}
		info, err := client.ManageGroupInfo(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
		req = &dcclient.ManageGroupInfoReq{
			Opt: def.OPT_MODIFY,
			Info: &dcclient.GroupInfo{
				GroupID: info.GroupID,
				Name        :  "测试组1修改",              //组名
				Uid        :    1699809227385606144,             //管理员用户id
			},
		}
		info, err = client.ManageGroupInfo(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dcclient.ManageGroupInfoReq{
			Opt: def.OPT_ADD,
			Info: &dcclient.GroupInfo{
				Name        :  "测试组2",              //组名
				Uid        :    1699809227385606145,             //管理员用户id
			},
		}
		info, err := client.ManageGroupInfo(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}