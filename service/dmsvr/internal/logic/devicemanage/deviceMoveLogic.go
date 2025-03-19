package devicemanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/topics"
	"gorm.io/gorm"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceMoveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceMoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceMoveLogic {
	return &DeviceMoveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceMoveLogic) DeviceMove(in *dm.DeviceMoveReq) (*dm.Empty, error) {
	if err := ctxs.IsAdmin(l.ctx); err != nil {
		return nil, err
	}
	l.ctx = ctxs.WithAllProject(l.ctx)
	diDB := relationDB.NewDeviceInfoRepo(l.ctx)
	oldDev, err := diDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.Old.ProductID, DeviceNames: []string{in.Old.DeviceName}})
	if err != nil {
		return nil, err
	}
	newDev, err := diDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.New.ProductID, DeviceNames: []string{in.New.DeviceName}})
	if err != nil {
		return nil, err
	}
	if oldDev.ProductID != newDev.ProductID && utils.SliceIn("profile", in.WithTarget...) {
		return nil, errors.Parameter.WithMsg("只有相同产品才可以迁移配置")
	}

	newDev.ProjectID = oldDev.ProjectID
	newDev.AreaID = oldDev.AreaID
	newDev.AreaIDPath = oldDev.AreaIDPath
	newDev.DeviceAlias = oldDev.DeviceAlias
	newDev.RatedPower = oldDev.RatedPower
	newDev.Address = oldDev.Address
	newDev.Tags = oldDev.Tags
	newDev.FirstBind = oldDev.FirstBind
	newDev.UserID = oldDev.UserID
	newDev.ExpTime = oldDev.ExpTime
	newDev.Distributor = oldDev.Distributor
	oldDev.ProjectID = def.NotClassified
	if oldDev.FirstBind.Valid && oldDev.FirstBind.Time.After(time.Now().AddDate(0, 0, -1)) { //绑定一天内的不算绑定时间
		oldDev.FirstBind.Valid = false
		oldDev.ExpTime.Valid = false
	} else {
		oldDev.ExpTime.Time = time.Now()
	}
	if newDev.ExpTime.Valid && newDev.ExpTime.Time.Before(time.Now()) { //如果过期了
		newDev.Status = def.DeviceStatusArrearage
	} else {
		newDev.Status = newDev.IsOnline + 1
	}
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, newDev.ProductID)
	if err != nil {
		return nil, err
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		diDB := relationDB.NewDeviceInfoRepo(tx)
		var updateProfile bool
		for _, target := range in.WithTarget {
			switch target {
			case "profile":
				updateProfile = true
				newDev.SchemaAlias = oldDev.SchemaAlias
			}
		}

		err := diDB.Update(l.ctx, oldDev)
		if err != nil {
			return err
		}
		err = diDB.Update(l.ctx, newDev)
		if err != nil {
			return err
		}
		{ //先删除新设备的关系
			err = relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
				ProductID:  newDev.ProductID,
				DeviceName: newDev.DeviceName,
			})
			if err != nil {
				return err
			}
			err = relationDB.NewDeviceProfileRepo(tx).DeleteByFilter(l.ctx,
				relationDB.DeviceProfileFilter{Device: devices.Core{
					ProductID: newDev.ProductID, DeviceName: newDev.DeviceName,
				}})
			if err != nil {
				return err
			}
			err = relationDB.NewUserDeviceCollectRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceCollectFilter{Cores: []*devices.Core{
				{ProductID: newDev.ProductID, DeviceName: newDev.DeviceName},
			}})
		}
		{ //将旧设备的关系迁移到新设备
			err = relationDB.NewUserDeviceShareRepo(tx).UpdateWithField(l.ctx, relationDB.UserDeviceShareFilter{
				ProductID:  oldDev.ProductID,
				DeviceName: oldDev.DeviceName,
			}, map[string]any{
				"product_id":  newDev.ProductID,
				"device_name": newDev.DeviceName,
			})
			if err != nil {
				return err
			}
			err = relationDB.NewUserDeviceCollectRepo(tx).UpdateWithField(l.ctx, relationDB.UserDeviceCollectFilter{Cores: []*devices.Core{
				{ProductID: oldDev.ProductID, DeviceName: oldDev.DeviceName},
			}}, map[string]any{
				"product_id":  newDev.ProductID,
				"device_name": newDev.DeviceName,
			})
			if updateProfile {
				err = relationDB.NewDeviceProfileRepo(tx).UpdateWithField(l.ctx,
					relationDB.DeviceProfileFilter{Device: devices.Core{ProductID: oldDev.ProductID, DeviceName: oldDev.DeviceName}}, map[string]any{"product_id": newDev.ProductID, "device_name": newDev.DeviceName})
				if err != nil {
					return err
				}
			}
		}
		if pi.DeviceType == def.DeviceTypeSubset { //如果是子设备需要将网关也切换过去
			gdDB := relationDB.NewGatewayDeviceRepo(tx)
			err = gdDB.DeleteDevAll(l.ctx, devices.Core{ProductID: newDev.ProductID, DeviceName: newDev.DeviceName})
			if err != nil {
				return err
			}
			err = gdDB.UpdateWithField(l.ctx, relationDB.GatewayDeviceFilter{SubDevice: &devices.Core{ProductID: oldDev.ProductID, DeviceName: oldDev.DeviceName}}, map[string]any{"product_id": newDev.ProductID, "device_name": newDev.DeviceName})
			if err != nil {
				return err
			}
		}
		if pi.DeviceType == def.DeviceTypeGateway {
			gdDB := relationDB.NewGatewayDeviceRepo(tx)
			err = gdDB.DeleteDevAll(l.ctx, devices.Core{ProductID: newDev.ProductID, DeviceName: newDev.DeviceName})
			if err != nil {
				return err
			}
			err = gdDB.UpdateWithField(l.ctx, relationDB.GatewayDeviceFilter{Gateway: &devices.Core{ProductID: oldDev.ProductID, DeviceName: oldDev.DeviceName}},
				map[string]any{"gateway_product_id": newDev.ProductID, "gateway_device_name": newDev.DeviceName})
			if err != nil {
				return err
			}
		}
		{
			err := l.svcCtx.AbnormalRepo.UpdateDevices(l.ctx, []*devices.Info{
				{ProductID: oldDev.ProductID, DeviceName: oldDev.DeviceName, TenantCode: string(oldDev.TenantCode),
					ProjectID: int64(oldDev.ProjectID), AreaID: int64(oldDev.AreaID), AreaIDPath: string(oldDev.AreaIDPath)},
				{ProductID: newDev.ProductID, DeviceName: newDev.DeviceName, TenantCode: string(newDev.TenantCode),
					ProjectID: int64(newDev.ProjectID), AreaID: int64(newDev.AreaID), AreaIDPath: string(newDev.AreaIDPath)}})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	oldDevCore := devices.Core{ProductID: oldDev.ProductID, DeviceName: oldDev.DeviceName}
	newDevCore := devices.Core{ProductID: newDev.ProductID, DeviceName: newDev.DeviceName}

	l.svcCtx.DeviceCache.SetData(l.ctx, oldDevCore, nil)
	l.svcCtx.DeviceCache.SetData(l.ctx, newDevCore, nil)
	err = DeleteDeviceTimeData(l.ctx, l.svcCtx, oldDev.ProductID, oldDev.DeviceName, DeleteModeAll)
	err = DeleteDeviceTimeData(l.ctx, l.svcCtx, newDev.ProductID, newDev.DeviceName, DeleteModeAll)
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmDeviceInfoUnbind, &oldDevCore)
	BindChange(l.ctx, l.svcCtx, pi, oldDevCore, int64(oldDev.ProjectID))
	BindChange(l.ctx, l.svcCtx, pi, newDevCore, int64(newDev.ProjectID))

	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, nil
}
