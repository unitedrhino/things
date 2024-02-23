package device

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dgsvr/client/deviceauth"
	"github.com/i-Things/things/service/dgsvr/pb/dg"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"sync"
)

func ThirdProtoLoginAuth(ctx context.Context, svcCtx *svc.ServiceContext, req *types.DeviceAuthLoginReq, cert []byte) error {
	pi, err := svcCtx.ProtocolM.ProtocolInfoIndex(ctx, &dm.ProtocolInfoIndexReq{TransProtocol: def.ProtocolMqtt})
	if err != nil {
		return err
	}
	var wait sync.WaitGroup
	var succ bool
	for _, v := range pi.List {
		wait.Add(1)
		go func(v *dm.ProtocolInfo) {
			defer wait.Done()
			var conf zrpc.RpcClientConf
			if v.EtcdKey != "" {
				conf.Etcd = svcCtx.Config.Etcd
				conf.Etcd.Key = v.EtcdKey
			} else if v.Endpoints != nil {
				conf.Endpoints = v.Endpoints
			} else { //如果都没有配置,那么就不走这个服务校验
				return
			}
			cli, err := zrpc.NewClient(conf)
			if err != nil {
				logx.WithContext(ctx).Errorf("NewClient ProtocolInfo:%#v err:%v", v, err)
				return
			}
			defer cli.Conn().Close()
			da := deviceauth.NewDeviceAuth(cli)
			_, err = da.LoginAuth(ctx, &dg.LoginAuthReq{Username: req.Username, //用户名
				Password:    req.Password, //密码
				ClientID:    req.ClientID, //clientID
				Ip:          req.Ip,       //访问的ip地址
				Certificate: cert,         //客户端证书
			})
			if err == nil {
				succ = true
			}
		}(v)
	}
	if succ {
		return nil
	}
	return errors.Permissions
}
func ThirdProtoAccessAuth(ctx context.Context, svcCtx *svc.ServiceContext, req *types.DeviceAuthAccessReq, action string) error {
	pi, err := svcCtx.ProtocolM.ProtocolInfoIndex(ctx, &dm.ProtocolInfoIndexReq{TransProtocol: def.ProtocolMqtt})
	if err != nil {
		return err
	}
	var wait sync.WaitGroup
	var succ bool
	for _, v := range pi.List {
		wait.Add(1)
		go func(v *dm.ProtocolInfo) {
			defer wait.Done()
			var conf zrpc.RpcClientConf
			if v.EtcdKey != "" {
				conf.Etcd = svcCtx.Config.Etcd
				conf.Etcd.Key = v.EtcdKey
			} else if v.Endpoints != nil {
				conf.Endpoints = v.Endpoints
			} else { //如果都没有配置,那么就不走这个服务校验
				return
			}
			cli, err := zrpc.NewClient(conf)
			if err != nil {
				logx.WithContext(ctx).Errorf("NewClient ProtocolInfo:%#v err:%v", v, err)
				return
			}
			defer cli.Conn().Close()
			da := deviceauth.NewDeviceAuth(cli)
			_, err = da.AccessAuth(ctx, &dg.AccessAuthReq{
				Username: req.Username,
				Topic:    req.Topic,
				ClientID: req.ClientID,
				Access:   action,
				Ip:       req.Ip,
			})
			if err == nil {
				logx.WithContext(ctx).Infof("AccessAuth ProtocolInfo:%#v succ", v)
				succ = true
			}
		}(v)
	}
	if succ {
		return nil
	}
	return errors.Permissions
}
