package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
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
		return nil
	})
	if err != nil {
		return nil, err
	}
	l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  oldDev.ProductID,
		DeviceName: oldDev.DeviceName,
	}, nil)
	l.svcCtx.DeviceCache.SetData(l.ctx, devices.Core{
		ProductID:  newDev.ProductID,
		DeviceName: newDev.DeviceName,
	}, nil)
	err = DeleteDeviceTimeData(l.ctx, l.svcCtx, oldDev.ProductID, oldDev.DeviceName)
	err = DeleteDeviceTimeData(l.ctx, l.svcCtx, newDev.ProductID, newDev.DeviceName)
	err = l.svcCtx.FastEvent.Publish(l.ctx, eventBus.DmDeviceInfoUnbind, &devices.Core{ProductID: oldDev.ProductID, DeviceName: oldDev.DeviceName})
	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, nil
}
