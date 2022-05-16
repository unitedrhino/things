package tests

import (
	"context"
	"fmt"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dcsvr/dc"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
	"testing"
)

func TestManageGroupMember(t *testing.T) {
	fmt.Println("TestManageGroupMember")
	client := dc.NewDc(zrpc.MustNewClient(zrpc.RpcClientConf{Etcd: discov.EtcdConf{
		Hosts: []string{"127.0.0.1:2379"},
		Key:   "dc.rpc",
	}}))
	ctx := context.Background()
	{
		req := &dc.ManageGroupMemberReq{
			Opt: def.OPT_ADD,
			Info: &dc.GroupMember{
				GroupID: 1710808183040118784, //组id
				//如果是用户,则是uid的十进制字符串,
				//如果是设备,则是productID:deviceName的组合方式
				MemberID:   "1699809227385606144", //成员id
				MemberType: 2,                     //成员类型:1:设备 2:用户
			},
		}
		info, err := client.ManageGroupMember(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
	{
		req := &dc.ManageGroupMemberReq{
			Opt: def.OPT_ADD,
			Info: &dc.GroupMember{
				GroupID: 1710808183040118784, //组id
				//如果是用户,则是uid的十进制字符串,
				//如果是设备,则是productID:deviceName的组合方式
				MemberID:   "21CYs1k9YpG:test8", //成员id
				MemberType: 1,                   //成员类型:1:设备 2:用户
			},
		}
		info, err := client.ManageGroupMember(ctx, req)
		t.Log(req, info)
		if err != nil {
			t.Errorf("%+v", errors.Fmt(err))
		}
	}
}
