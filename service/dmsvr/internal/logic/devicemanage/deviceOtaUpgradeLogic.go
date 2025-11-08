package devicemanagelogic

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	otamanagelogic "gitee.com/unitedrhino/things/service/dmsvr/internal/logic/otamanage"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgOta"
	"gitee.com/unitedrhino/things/share/userSubscribe"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceOtaUpgradeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceOtaUpgradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceOtaUpgradeLogic {
	return &DeviceOtaUpgradeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备能升级的升级包
func (l *DeviceOtaUpgradeLogic) DeviceOtaUpgrade(in *dm.DeviceOtaUpgradeReq) (*dm.DeviceOtaUpgradeResp, error) {
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	})
	if err != nil {
		return nil, err
	}

	if di.Version.GetValue() != in.Version {
		//如果不一样则需要判断是否是ota升级的,如果是,则需要更新升级状态
		dfs, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
			ProductID:   di.ProductID,
			DeviceNames: []string{di.DeviceName},
			DestVersion: di.Version.GetValue(),
			Statues: []int64{msgOta.DeviceStatusConfirm, msgOta.DeviceStatusQueued,
				msgOta.DeviceStatusNotified, msgOta.DeviceStatusInProgress,
				msgOta.DeviceStatusCanceled, msgOta.DeviceStatusFailure}, //除了成功的都过滤出来
		}, nil)
		if err != nil {
			if !errors.Cmp(err, errors.NotFind) {
				return nil, err
			}
		} else {
			var once sync.Once
			for _, df := range dfs {
				df.Step = 100
				df.Status = msgOta.DeviceStatusSuccess
				df.Detail = "升级成功"
				err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, df)
				if err != nil {
					l.Error(err)
					continue
				}
				once.Do(func() {
					appMsg := application.OtaReport{
						Device:    devices.Core{ProductID: di.ProductID, DeviceName: di.DeviceName},
						Timestamp: time.Now().UnixMilli(), Status: df.Status, Detail: df.Detail, Step: df.Step,
					}
					err = l.svcCtx.UserSubscribe.Publish(l.ctx, userSubscribe.DeviceOtaReport, appMsg, map[string]any{
						"productID":  di.ProductID,
						"deviceName": di.DeviceName,
					}, map[string]any{
						"projectID": di.ProjectID,
					}, map[string]any{
						"projectID": cast.ToString(di.ProjectID),
						"areaID":    cast.ToString(di.AreaID),
					})
				})
			}
			err := relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID, DeviceName: in.DeviceName}, map[string]any{
				"version":              in.Version,
				"need_confirm_job_id":  "",
				"need_confirm_version": 0,
			})
			if err != nil {
				return nil, err
			}
			return &dm.DeviceOtaUpgradeResp{}, nil
		}
		return &dm.DeviceOtaUpgradeResp{}, nil
	}

	df, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
		ProductID:    in.ProductID,
		DeviceNames:  []string{in.DeviceName},
		WithFirmware: true,
		WithJob:      true,
		DestVersion:  in.Version,
		Statues:      []int64{msgOta.DeviceStatusQueued},
	})
	if err != nil && !errors.Cmp(err, errors.NotFind) {
		return nil, err
	}

	if df == nil {
		return nil, errors.NotFind
	}
	data, err := otamanagelogic.GenUpgradeParams(l.ctx, l.svcCtx, df.Firmware, df.Files)
	if err != nil {
		return nil, err
	}
	var resp = dm.DeviceOtaUpgradeResp{
		Firmware: &dm.OtaUpgradeData{
			Version:    data.Version,
			Name:       df.Firmware.Name,
			Desc:       df.Firmware.Desc,
			IsDiff:     data.IsDiff,
			SignMethod: data.SignMethod,
			Extra:      data.Extra,
		},
	}
	if len(data.Files) == 0 {
		resp.Firmware.Files = append(resp.Firmware.Files, utils.Copy[dm.OtaFile](data.File))
	} else {
		resp.Firmware.Files = utils.CopySlice[dm.OtaFile](data.Files)
	}
	if in.StartUpdate {
		df.Status = msgOta.DeviceStatusNotified
		df.Detail = "接口获取升级"
		df.PushTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		err = relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, df)
		if err != nil {
			logx.WithContext(l.ctx).Error(err)
			return nil, err
		}
	}
	return &resp, nil
}
