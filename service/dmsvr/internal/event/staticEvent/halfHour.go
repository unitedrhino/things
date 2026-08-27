package staticEvent

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	cacheRepo "gitee.com/unitedrhino/things/service/dmsvr/internal/repo/cache"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/tsDB/schemaDataRepo"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"sort"
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
	w.Add(4)
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.ProductCategoryStatic()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.TimescaleHandle()
		if err != nil {
			l.Error(err)
		}
	})
	utils.Go(l.ctx, func() {
		defer w.Done()
		err := l.runDeviceStatusMaintenance()
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
	w.Wait()
	return nil
}

func (l *HalfHourHandle) runDeviceStatusMaintenance() error {
	lock := cacheRepo.NewDeviceStatusMaintenanceLock(l.svcCtx.Cache)
	token, acquired, err := lock.TryLock(l.ctx)
	if err != nil {
		return fmt.Errorf("acquire device status maintenance lock: %w", err)
	}
	if !acquired {
		l.Info("skip duplicate device status maintenance")
		return nil
	}
	defer func() {
		if err := lock.Unlock(l.ctx, token); err != nil {
			l.Errorf("release device status maintenance lock: %v", err)
		}
	}()

	return runDeviceStatusStages(l.DeviceExp, l.DeviceAbnormalRecover, l.DeviceAbnormalSet)
}

func runDeviceStatusStages(expire, recover, markAbnormal func() error) error {
	if err := expire(); err != nil {
		return fmt.Errorf("update expired device status: %w", err)
	}
	if err := recover(); err != nil {
		return fmt.Errorf("recover abnormal device status: %w", err)
	}
	if err := markAbnormal(); err != nil {
		return fmt.Errorf("mark abnormal device status: %w", err)
	}
	return nil
}

// 参考: https://docs.tigerdata.com/api/latest/continuous-aggregates/refresh_continuous_aggregate/
func (l *HalfHourHandle) TimescaleHandle() error { //timescale 视图更新
	if stores.GetTsDBType() != conf.Pgsql {
		return nil
	}
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var template = fmt.Sprintf("CALL refresh_continuous_aggregate('%%s', '%s', '%s', force => TRUE);", yesterday, tomorrow)
	for _, tbName := range schemaDataRepo.TableNames {
		stores.GetTsConn(l.ctx).Exec(fmt.Sprintf(template, tbName+"_day"))
		stores.GetTsConn(l.ctx).Exec(fmt.Sprintf(template, tbName+"_hour"))
	}
	return nil
}

func (l *HalfHourHandle) DeviceExp() error { //设备过期处理
	{ //有效期到了之后不启用
		start := time.Now()
		updated, batches, err := relationDB.NewDeviceInfoRepo(l.ctx).UpdateExpiredDeviceStatus(l.ctx, start)
		if err != nil {
			return err
		}
		l.Infof("expired device status maintenance updated=%d batches=%d duration=%s", updated, batches, time.Since(start))
	}
	{ //清除设置了过期时间且过期了的分享
		err := relationDB.NewUserDeviceShareRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
			ExpTime: stores.CmpAnd(stores.CmpLte(time.Now()), stores.CmpIsNull(false)),
		})
		if err != nil {
			return err
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
	var recoverDeviceIDs []int64
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
		recoverDeviceIDs = append(recoverDeviceIDs, d.ID)
	}
	if len(recoverDeviceIDs) > 0 {
		sort.Slice(recoverDeviceIDs, func(i, j int) bool { return recoverDeviceIDs[i] < recoverDeviceIDs[j] })
		repo := relationDB.NewDeviceInfoRepo(l.ctx)
		for start := 0; start < len(recoverDeviceIDs); start += relationDB.DeviceStatusUpdateBatchSize {
			end := min(start+relationDB.DeviceStatusUpdateBatchSize, len(recoverDeviceIDs))
			updated, err := repo.RecoverAbnormalDeviceStatusBatch(l.ctx, recoverDeviceIDs[start:end])
			if err != nil {
				return err
			}
			l.Infof("recover abnormal devices updated=%d", len(updated))
			for _, v := range updated {
				if err := l.svcCtx.AbnormalRepo.Insert(l.ctx, &deviceLog.Abnormal{
					TenantCode: string(v.TenantCode),
					ProjectID:  int64(v.ProjectID),
					AreaID:     int64(v.AreaID),
					AreaIDPath: string(v.AreaIDPath),
					ProductID:  v.ProductID,
					DeviceName: v.DeviceName,
					Action:     def.False,
					Type:       "online", //上下线异常
					Timestamp:  time.Now(),
					Reason:     "设备异常上下线恢复",
				}); err != nil {
					return err
				}
			}
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
	var abnormalDeviceIDs []int64
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
		abnormalDeviceIDs = append(abnormalDeviceIDs, d.ID)
	}
	if len(abnormalDeviceIDs) > 0 {
		sort.Slice(abnormalDeviceIDs, func(i, j int) bool { return abnormalDeviceIDs[i] < abnormalDeviceIDs[j] })
		repo := relationDB.NewDeviceInfoRepo(l.ctx)
		for start := 0; start < len(abnormalDeviceIDs); start += relationDB.DeviceStatusUpdateBatchSize {
			end := min(start+relationDB.DeviceStatusUpdateBatchSize, len(abnormalDeviceIDs))
			updated, err := repo.MarkAbnormalDeviceStatusBatch(l.ctx, abnormalDeviceIDs[start:end])
			if err != nil {
				return err
			}
			l.Infof("mark abnormal devices updated=%d", len(updated))
			for _, v := range updated {
				if err := l.svcCtx.AbnormalRepo.Insert(l.ctx, &deviceLog.Abnormal{
					TenantCode: string(v.TenantCode),
					ProjectID:  int64(v.ProjectID),
					AreaID:     int64(v.AreaID),
					AreaIDPath: string(v.AreaIDPath),
					ProductID:  v.ProductID,
					DeviceName: v.DeviceName,
					Action:     def.True,
					Type:       "online", //上下线异常
					Timestamp:  time.Now(),
					Reason:     "设备异常频繁上下线",
				}); err != nil {
					return err
				}
			}
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
