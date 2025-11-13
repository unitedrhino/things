package staticEvent

import (
	"context"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/sdk/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/zeromicro/go-zero/core/logx"
)

type OneHourHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewOneHourHandle(ctx context.Context, svcCtx *svc.ServiceContext) *OneHourHandle {
	return &OneHourHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *OneHourHandle) Handle() error { //产品品类设备数量统计
	err := l.DeviceOnlineFix()
	if err != nil {
		l.Error(err)
	}
	return nil
}

func (l *OneHourHandle) DeviceOnlineFix() error { //设备在线修复
	devs, err := relationDB.NewDeviceInfoRepo(l.ctx).FindCoreByFilter(l.ctx, relationDB.DeviceFilter{IsOnline: def.True}, nil)
	if err != nil {
		return err
	}
	devMap, err := protocol.GetActivityDevices(l.ctx)
	if err != nil {
		l.Error(err)
		return err
	}
	var needOnline []devices.Core
	for _, d := range devs {
		if _, ok := devMap[d]; ok {
			continue
		}
		//如果线上没有,但是这里有,需要进行处理
		needOnline = append(needOnline, d)
	}
	if len(needOnline) > 0 {
		l.Infof("DeviceOnlineFix.UpdatesDeviceActivity devs:%v", utils.Fmt(needOnline))
		protocol.UpdatesDeviceActivity(l.ctx, needOnline)
	}
	return nil
}
