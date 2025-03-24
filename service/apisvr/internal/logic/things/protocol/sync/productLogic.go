package sync

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/protocolSyncClient"
	"github.com/zeromicro/go-zero/zrpc"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 产品同步(如果该协议不支持会返回不支持)
func NewProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductLogic {
	return &ProductLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductLogic) Product(req *types.CodeReq) error {
	pi, err := l.svcCtx.ProtocolM.ProtocolInfoRead(l.ctx, &dm.WithIDCode{Code: req.Code})
	if err != nil {
		return err
	}
	if pi.IsEnableSyncProduct != def.True {
		return errors.NotRealize.WithMsg("该协议不支持设备同步")
	}
	var conf = l.svcCtx.Config.DgRpc.Conf
	if pi.EtcdKey != "" {
		conf.Etcd = l.svcCtx.Config.Etcd
		conf.Etcd.Key = pi.EtcdKey
	} else if pi.Endpoints != nil {
		conf.Endpoints = pi.Endpoints
	} else { //如果都没有配置,那么就不走这个服务校验
		return errors.System.AddMsg("连接不到协议网关").AddDetail(err)

	}
	cli, err := zrpc.NewClient(conf)
	if err != nil {
		return errors.System.AddMsg("连接不到协议网关").AddDetail(err)
	}
	_, err = protocolSyncClient.NewProtocolSync(cli).SyncProduct(l.ctx, &protocolSyncClient.Empty{})
	if err != nil {
		return err
	}
	return nil
}
