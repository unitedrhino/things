package devicemanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gorm.io/gorm"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceResetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceResetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceResetLogic {
	return &DeviceResetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeviceResetLogic) DeviceReset(in *dm.DeviceResetReq) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	dev := devices.Core{
		ProductID:  in.Device.ProductID,
		DeviceName: in.Device.DeviceName,
	}
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, dev)
	if err != nil {
		return nil, err
	}
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, di.ProductID)
	if err == nil {
		return nil, err
	}
	if in.Log {
		err := DeleteDeviceTimeData(l.ctx, l.svcCtx, di.ProductID, di.DeviceName, DeleteModeAll)
		if err != nil {
			return nil, err
		}
	} else if in.DeviceSchema {
		err = l.svcCtx.SchemaManaRepo.DeleteDeviceProperty(l.ctx, dev.ProductID, dev.DeviceName, nil)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("%s.SchemaManaRepo.DeleteDevice err=%v", utils.FuncName(), err)
			return nil, err
		}
	}
	var updateMap = make(map[string]interface{})
	areaID := dataType.AreaID(def.NotClassified)
	var projectID dataType.ProjectID
	var areaIDPath string
	if in.Bind {
		areaID = dataType.AreaID(def.NotClassified)
		ti, err := l.svcCtx.TenantCache.GetData(l.ctx, def.TenantCodeDefault)
		if err != nil {
			return nil, err
		}
		projectID = dataType.ProjectID(ti.DefaultProjectID)
		if ti.DefaultAreaID != 0 {
			areaID = dataType.AreaID(ti.DefaultAreaID)
		}

		ai, err := l.svcCtx.AreaCache.GetData(l.ctx, int64(areaID))
		if err != nil {
			return nil, err
		}
		areaIDPath = ai.AreaIDPath
		updateMap["project_id"] = projectID
		updateMap["area_id"] = areaID
		updateMap["area_id_path"] = areaIDPath
		updateMap["user_id"] = def.RootNode
	}
	if in.Info {
		updateMap["device_alias"] = fmt.Sprintf("%s%d", pi.ProductName, GenID())
		updateMap["position"] = nil
		updateMap["rated_power"] = 0
		updateMap["imei"] = ""
		updateMap["mac"] = ""
		updateMap["version"] = ""
		updateMap["hard_info"] = ""
		updateMap["soft_info"] = ""
		updateMap["mobile_operator"] = ""
		updateMap["phone"] = ""
		updateMap["iccid"] = ""
		updateMap["address"] = ""
		updateMap["adcode"] = ""
		updateMap["tags"] = map[string]string{}
		updateMap["schema_alias"] = map[string]string{}
		updateMap["rssi"] = ""
		updateMap["protocol_conf"] = map[string]string{}
		updateMap["first_login"] = nil
		updateMap["first_bind"] = nil
		updateMap["last_login"] = nil
		updateMap["log_level"] = 1
		updateMap["status"] = di.IsOnline + 1
		updateMap["exp_time"] = nil
		updateMap["desc"] = ""
	}
	err = stores.GetTenantConn(l.ctx).Transaction(func(tx *gorm.DB) error {
		if len(updateMap) != 0 {
			err := relationDB.NewDeviceInfoRepo(l.ctx).UpdateWithField(l.ctx, relationDB.DeviceFilter{ProductID: di.ProductID, DeviceNames: []string{di.DeviceName}}, updateMap)
			if err != nil {
				return nil
			}
		}
		if in.Info {
			err = relationDB.NewDeviceProfileRepo(tx).DeleteByFilter(l.ctx, relationDB.DeviceProfileFilter{Device: dev})
			if err != nil {
				l.Errorf("%s.NewDeviceProfileRepo.Delete err=%+v", utils.FuncName(), err)
				return err
			}
			err = relationDB.NewUserDeviceShareRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceShareFilter{
				ProductID:  di.ProductID,
				DeviceName: di.DeviceName,
			})
			if err != nil {
				return err
			}
			err = relationDB.NewUserDeviceCollectRepo(tx).DeleteByFilter(l.ctx, relationDB.UserDeviceCollectFilter{Cores: []*devices.Core{
				{ProductID: di.ProductID, DeviceName: di.DeviceName},
			}})
			if err != nil {
				return err
			}
			err = relationDB.NewGatewayDeviceRepo(tx).DeleteDevAll(l.ctx, dev)
			if err != nil {
				return err
			}
		}
		if in.DeviceSchema {
			err = relationDB.NewDeviceSchemaRepo(tx).DeleteByFilter(l.ctx, relationDB.DeviceSchemaFilter{ProductID: di.ProductID, DeviceName: di.DeviceName})
			if err != nil {
				return err
			}
		}
		if in.Bind {
			err := l.svcCtx.AbnormalRepo.UpdateDevices(l.ctx, []*devices.Info{
				{ProductID: di.ProductID, DeviceName: di.DeviceName,
					ProjectID: int64(projectID), AreaID: int64(areaID), AreaIDPath: areaIDPath}})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return &dm.Empty{}, err
}
