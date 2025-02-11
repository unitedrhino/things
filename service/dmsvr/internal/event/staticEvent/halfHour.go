package staticEvent

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/share/devices"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type HalfHourHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewHalfHourHandle(ctx context.Context, svcCtx *svc.ServiceContext) *HalfHourHandle {
	return &HalfHourHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *HalfHourHandle) Handle() error { //产品品类设备数量统计
	w := sync.WaitGroup{}
	w.Add(5)
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.ProductCategoryStatic()
		if err != nil {
			l.Error(err)
		}
	})

	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.DeviceExp()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.DeviceMsgCount()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.DeviceAbnormalRecover()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.DeviceAbnormalSet()
		if err != nil {
			l.Error(err)
		}
	})
	w.Wait()
	return nil
}

func (l *HalfHourHandle) DeviceExp() error { //设备过期处理
	{ //有效期到了之后不启用
		err := relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx,
			relationDB.DeviceFilter{HasOwner: def.True, ExpTime: stores.CmpAnd(stores.CmpLte(time.Now()), stores.CmpIsNull(false))},
			map[string]any{"status": def.DeviceStatusArrearage})
		if err != nil {
			l.Error(err)
		}
	}
	{ //清除设置了过期时间且过期了的分享
		err := relationDB.NewUserDeviceShareRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
			ExpTime: stores.CmpAnd(stores.CmpLte(time.Now()), stores.CmpIsNull(false)),
		})
		if err != nil {
			l.Error(err)
		}
	}
	return nil
}
func (l *HalfHourHandle) DeviceAbnormalRecover() error { //设备上下线异常恢复
	now := time.Now()
	dis, err := relationDB.NewDeviceInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceFilter{
		Statuses: []int64{def.DeviceStatusAbnormal},
	}, nil)
	if err != nil {
		return err
	}
	var recoverDevices []*devices.Core
	var recoverDeviceDetail []*relationDB.DmDeviceInfo
	for _, d := range dis {
		count, err := l.svcCtx.StatusRepo.GetCountLog(l.ctx, deviceLog.StatusFilter{
			ProductID:  d.ProductID,
			DeviceName: d.DeviceName,
		}, def.PageInfo2{
			TimeStart: now.Add(-time.Minute * 60).UnixMilli(),
		})
		if err != nil {
			continue
		}
		if count > 5 { //如果前一个小时还超过5次的登入登出,则保持异常状态
			continue
		}
		recoverDeviceDetail = append(recoverDeviceDetail, d)
		recoverDevices = append(recoverDevices, &devices.Core{
			ProductID:  d.ProductID,
			DeviceName: d.DeviceName,
		})
	}
	if len(recoverDeviceDetail) > 0 {
		l.Infof("recoverDevices:%v", utils.Fmt(recoverDevices))
		err := relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx,
			relationDB.DeviceFilter{Cores: recoverDevices},
			map[string]any{"status": stores.Expr("is_online + 1")})
		if err != nil {
			l.Error(err)
		}
		for _, v := range recoverDeviceDetail {
			l.svcCtx.AbnormalRepo.Insert(l.ctx, &deviceLog.Abnormal{
				TenantCode: v.TenantCode,
				ProjectID:  v.ProjectID,
				AreaID:     v.AreaID,
				AreaIDPath: v.AreaIDPath,
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				Action:     def.False,
				Type:       "online", //上下线异常
				Timestamp:  time.Now(),
				Reason:     "设备异常上下线恢复",
			})
		}
	}
	return nil
}

func (l *HalfHourHandle) DeviceAbnormalSet() error { //设备上下线异常设置
	now := time.Now()
	dis, err := relationDB.NewDeviceInfoRepo(l.ctx).FindByFilter(l.ctx, relationDB.DeviceFilter{
		LastLoginTime: &def.TimeRange{
			Start: now.Add(-time.Minute * 60).Unix(),
		},
		Statuses: []int64{def.DeviceStatusOnline, def.DeviceStatusOffline},
	}, nil)
	if err != nil {
		return err
	}
	var abnormalDevices []*devices.Core
	var abnormalDeviceDetail []*relationDB.DmDeviceInfo
	for _, d := range dis {
		count, err := l.svcCtx.StatusRepo.GetCountLog(l.ctx, deviceLog.StatusFilter{
			ProductID:  d.ProductID,
			DeviceName: d.DeviceName,
		}, def.PageInfo2{
			TimeStart: now.Add(-time.Minute * 60).UnixMilli(),
		})
		if err != nil {
			continue
		}
		if count < 10 {
			continue
		}
		//如果一个小时内上下线次数大于10次,则判断为异常设备
		abnormalDeviceDetail = append(abnormalDeviceDetail, d)
		abnormalDevices = append(abnormalDevices, &devices.Core{
			ProductID:  d.ProductID,
			DeviceName: d.DeviceName,
		})
	}
	if len(abnormalDeviceDetail) > 0 {
		l.Infof("abnormalDevices:%v", utils.Fmt(abnormalDevices))
		err := relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx,
			relationDB.DeviceFilter{Cores: abnormalDevices, Statuses: []int64{def.DeviceStatusOnline, def.DeviceStatusOffline}},
			map[string]any{"status": def.DeviceStatusAbnormal})
		if err != nil {
			l.Error(err)
		}
		for _, v := range abnormalDeviceDetail {
			l.svcCtx.AbnormalRepo.Insert(l.ctx, &deviceLog.Abnormal{
				TenantCode: v.TenantCode,
				ProjectID:  v.ProjectID,
				AreaID:     v.AreaID,
				AreaIDPath: v.AreaIDPath,
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
				Action:     def.True,
				Type:       "online", //上下线异常
				Timestamp:  time.Now(),
				Reason:     "设备异常频繁上下线",
			})
		}
	}
	return nil
}

func (l *HalfHourHandle) DeviceMsgCount() error { //产品品类设备数量统计
	end := time.Now()
	var fm = end.Minute() / 30 * 30
	var countData []*relationDB.DmDeviceMsgCount
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), fm, 0, 0, time.Local)
	start := end.Add(-time.Minute * 30)
	{
		hubCount, err := l.svcCtx.HubLogRepo.GetCountLog(l.ctx, deviceLog.HubFilter{}, def.PageInfo2{
			TimeStart: start.UnixMilli(),
			TimeEnd:   end.UnixMilli(),
		})
		if err != nil {
			l.Error(err)
		}
		countData = append(countData, &relationDB.DmDeviceMsgCount{
			Type: deviceLog.MsgTypePublish,
			Num:  hubCount,
			Date: end,
		})
	}
	{
		sendCount, err := l.svcCtx.SendRepo.GetCountLog(l.ctx, deviceLog.SendFilter{}, def.PageInfo2{
			TimeStart: start.UnixMilli(),
			TimeEnd:   end.UnixMilli(),
		})
		if err != nil {
			l.Error(err)
		}
		countData = append(countData, &relationDB.DmDeviceMsgCount{
			Type: deviceLog.MsgTypeSend,
			Num:  sendCount,
			Date: end,
		})
	}
	{
		t, err := relationDB.NewDeviceInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeviceFilter{IsOnline: def.True})
		if err != nil {
			l.Error(err)
		}
		countData = append(countData, &relationDB.DmDeviceMsgCount{
			Type: deviceLog.MsgTypeOnline,
			Num:  t,
			Date: end,
		})
	}
	err := relationDB.NewDeviceMsgCountRepo(l.ctx).MultiInsert(l.ctx, countData)
	if err != nil {
		l.Error(err)
	}
	return nil
}

func (l *HalfHourHandle) ProductCategoryStatic() error { //产品品类设备数量统计
	pcDB := relationDB.NewProductCategoryRepo(l.ctx)
	pcs, err := pcDB.FindByFilter(l.ctx, relationDB.ProductCategoryFilter{}, nil)
	if err != nil {
		return err
	}
	for _, pc := range pcs {
		ids := utils.GetIDPath(pc.IDPath)
		total, err := relationDB.NewDeviceInfoRepo(l.ctx).CountByFilter(l.ctx, relationDB.DeviceFilter{ProductCategoryIDs: ids})
		if err != nil {
			l.Error(err)
			continue
		}
		pc.DeviceCount = total
		err = pcDB.Update(l.ctx, pc)
		if err != nil {
			l.Error(err)
			continue
		}
	}
	return nil
}
