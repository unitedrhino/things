package staticEvent

import (
	"context"
	"time"

	"gitee.com/unitedrhino/core/service/syssvr/pb/sys"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/sdk/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
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
	err = l.DeviceStatic()
	if err != nil {
		l.Error(err)
	}
	return nil
}

func (l *OneHourHandle) DeviceStatic() error { //区域下的设备数量统计
	{
		ret, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, &sys.ProjectInfoIndexReq{})
		if err != nil {
			return err
		}
		var projectIDs []int64
		for _, v := range ret.List {
			projectIDs = append(projectIDs, v.ProjectID)
		}
		err = logic.FillProjectDeviceCount(l.ctx, l.svcCtx, projectIDs...)
		time.Sleep(time.Second * 30) //休息一下减少波峰
	}
	{
		var total int64 = 9999 //如果三次都没有成功自然退出
		var size int64 = 500
		var areas []*sys.AreaInfo
		var errCount int64 = 0
		for page := int64(0); page*size < total; page++ {
			err := func() error {
				ret, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{Page: &sys.PageInfo{
					Page: page + 1,
					Size: size,
				}})
				if err != nil {
					return err
				}
				total = ret.Total
				for _, v := range ret.List {
					areas = append(areas, v)
				}
				return nil
			}()
			if err != nil {
				l.Error(err)
				errCount++
			}
			if errCount > 3 { //只有三次错误的机会
				break
			}
		}
		err := logic.DirectFillAreaDeviceCount(l.ctx, l.svcCtx, areas...)
		if err != nil {
			l.Error(err)
		}
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
