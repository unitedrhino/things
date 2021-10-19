package tests

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

func TestManageGroupMember(t *testing.T) {
	fmt.Println("TestManageGroupMember")
	client := dcclient.NewDc(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dc.rpc",
	}}))
	ctx := context.Background()
	{
		req := &dcclient.ManageGroupMemberReq{
			Opt: def.OPT_ADD,
			Info: &dcclient.GroupMember{
				GroupID :1710808183040118784, //组id
				//如果是用户,则是uid的十进制字符串,
				//如果是设备,则是productID:deviceName的组合方式
				MemberID    :"1699809227385606144",        //成员id
				MemberType  :2,   //成员类型:1:设备 2:用户
			},
		}
		info, err := client.ManageGroupMember(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dcclient.ManageGroupMemberReq{
			Opt: def.OPT_ADD,
			Info: &dcclient.GroupMember{
				GroupID :1710808183040118784, //组id
				//如果是用户,则是uid的十进制字符串,
				//如果是设备,则是productID:deviceName的组合方式
				MemberID    :"21CYs1k9YpG:test8",        //成员id
				MemberType  :1,   //成员类型:1:设备 2:用户
			},
		}
		info, err := client.ManageGroupMember(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}