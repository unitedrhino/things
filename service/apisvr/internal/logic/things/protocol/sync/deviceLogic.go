package sync

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/protocolSyncClient"
	"github.com/zeromicro/go-zero/zrpc"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 设备同步(如果该协议不支持会返回不支持)
func NewDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceLogic {
	return &DeviceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeviceLogic) Device(req *types.ProtocolSyncDeviceReq) (resp *types.ProtocolSyncDeviceResp, err error) {
	pi, err := l.svcCtx.ProtocolM.ProtocolInfoRead(l.ctx, &dm.WithIDCode{Code: req.Code})
	if err != nil {
		return nil, err
	}
	if pi.IsEnableSyncDevice != def.True {
		return nil, errors.NotRealize.WithMsg("该协议不支持设备同步")
	}
	var conf = l.svcCtx.Config.DgRpc.Conf
	if pi.EtcdKey != "" {
		conf.Etcd = l.svcCtx.Config.Etcd
		conf.Etcd.Key = pi.EtcdKey
	} else if pi.Endpoints != nil {
		conf.Endpoints = pi.Endpoints
	} else { //如果都没有配置,那么就不走这个服务校验
		return nil, errors.System.AddMsg("连接不到协议网关").AddDetail(err)
	}
	cli, err := zrpc.NewClient(conf)
	if err != nil {
		return nil, errors.System.AddMsg("连接不到协议网关").AddDetail(err)
	}
	ret, err := protocolSyncClient.NewProtocolSync(cli).SyncDevice(l.ctx, &protocolSyncClient.SyncDeviceReq{ProductID: req.ProductID})
	if err != nil {
		return nil, err
	}
	return utils.Copy[types.ProtocolSyncDeviceResp](ret), nil
}
