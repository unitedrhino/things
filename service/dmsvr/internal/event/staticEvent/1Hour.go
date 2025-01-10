package staticEvent

import (
	"context"
	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/sdk/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
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
	w := sync.WaitGroup{}
	w.Add(2)
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.DeviceOnlineFix()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.DeviceStatic()
		if err != nil {
			l.Error(err)
		}
	})
	w.Wait()
	return nil
}
func (l *OneHourHandle) DeviceStatic() error { //区域下的设备数量统计
	err := func() error {
		ret, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{})
		if err != nil {
			return err
		}
		var areaPaths []string
		for _, v := range ret.List {
			areaPaths = append(areaPaths, v.AreaIDPath)
		}
		err = logic.FillAreaDeviceCount(l.ctx, l.svcCtx, areaPaths...)
		return err
	}()
	if err != nil {
		l.Error(err)
	}
	ret, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, &sys.ProjectInfoIndexReq{})
	if err != nil {
		return err
	}
	var projectIDs []int64
	for _, v := range ret.List {
		projectIDs = append(projectIDs, v.ProjectID)
	}
	err = logic.FillProjectDeviceCount(l.ctx, l.svcCtx, projectIDs...)
	return err
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
