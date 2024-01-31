package otataskmanagelogic

import (
	"context"
	"encoding/json"
	"errors"
	"gitee.com/i-Things/core/shared/utils"

	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"gitee.com/i-Things/core/shared/def"
	thingsError "gitee.com/i-Things/core/shared/errors"
	"github.com/hashicorp/go-uuid"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB *relationDB.DeviceInfoRepo
	OfDB *relationDB.OtaFirmwareRepo
}

func NewOtaTaskCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskCreateLogic {
	return &OtaTaskCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
	}
}
func (l *OtaTaskCreateLogic) CheckTask(in *dm.OtaTaskCreateReq) (bool, error) {
	/**
	upgradeType=1 静态升级，批量升级当前已有设备
	type=1 全部设备   可指定待升级版本号，不指定则是全部版本
	type=2 定向升级   必须勾选设备
		**/
	/**upgradeType=2 动态升级，持续升级新增设备
	type=1 全部设备   必须指定待升级版本号
	    没有定向升级
	**/
	var params []string
	if in.UpgradeType == 1 {
		if in.Type == 1 {
			in.DeviceList = nil
			if in.VersionList != nil {
				err := json.Unmarshal([]byte(in.VersionList.Value), &params)
				if err != nil {
					return false, errors.New("版本信息错误" + err.Error())
				}
			}
		}
		if in.Type == 2 {
			in.VersionList = nil
			if in.DeviceList == nil {
				return false, errors.New("定向升级需指定设备")
			}
		}
	}
	if in.UpgradeType == 2 {
		if in.Type == 2 {
			return false, errors.New("动态升级不能定向升级")
		}
		if in.VersionList == nil {
			return false, errors.New("动态升级需指定版本")
		}
		err := json.Unmarshal([]byte(in.VersionList.Value), &params)
		if err != nil {
			return false, errors.New("版本信息错误" + err.Error())
		}
	}
	_, err := l.OfDB.FindOne(l.ctx, in.FirmwareID)
	if thingsError.Cmp(err, thingsError.NotFind) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// 创建批量升级任务
func (l *OtaTaskCreateLogic) OtaTaskCreate(in *dm.OtaTaskCreateReq) (*dm.OtaTaskCreatResp, error) {
	productFirmwareInfo, err := l.OfDB.FindOne(l.ctx, in.FirmwareID)
	if thingsError.Cmp(err, thingsError.NotFind) {
		return nil, errors.New("未知固件信息")
	} else if err != nil {
		return nil, err
	}

	find, err := l.CheckTask(in)
	if err != nil {
		l.Errorf("AddOtaTask|CheckFirmware|in=%v\n", in)
		return nil, thingsError.Database.AddDetail(err)
	} else if find == false {
		return nil, thingsError.Parameter.AddDetail("not find firmware id:" + utils.ToString(in.FirmwareID))
	}
	taskUid, _ := uuid.GenerateUUID()
	di := relationDB.DmOtaTask{
		ProductID:   productFirmwareInfo.ProductID,
		FirmwareID:  in.FirmwareID,
		TaskUid:     taskUid,
		Type:        int64(in.Type),
		UpgradeType: int64(in.UpgradeType),
		AutoRepeat:  int64(in.AutoRepeat),
		Status:      1,
	}
	var deviceList []string
	if in.DeviceList != nil {
		err = json.Unmarshal([]byte(in.DeviceList.Value), &deviceList)
		if err != nil {
			return nil, thingsError.Parameter.AddDetail("deviceList need json")
		}
		di.DeviceList = in.DeviceList.Value
	} else {
		di.DeviceList = "{}"
	}
	var versionList []string
	if in.VersionList != nil {
		err = json.Unmarshal([]byte(in.VersionList.Value), &versionList)
		if err != nil {
			return nil, thingsError.Parameter.AddDetail("VersionList need json")
		}
		di.VersionList = in.VersionList.Value
	} else {
		di.VersionList = "{}"
	}
	err = relationDB.NewOtaTaskRepo(l.ctx).Insert(l.ctx, &di)
	if err != nil {
		l.Errorf("AddOtaTask.OtaTaskInfo.Insert err=%+v", err)
		return nil, thingsError.System.AddDetail(err)
	}
	taskID := di.ID
	//插入taskDevice
	var otDB = relationDB.NewOtaTaskDevicesRepo(l.ctx)
	if in.Type == 1 {
		di, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{ProductID: di.ProductID}, &def.PageInfo{Size: 100, Page: 1})
		if err != nil {
			return nil, err
		}

		for _, v := range di {
			otDB.Insert(l.ctx, &relationDB.DmOtaTaskDevices{
				FirmwareID:    in.FirmwareID,
				TaskUid:       taskUid,
				ProductID:     productFirmwareInfo.ProductID,
				Status:        101,
				TargetVersion: v.Version,
				DeviceName:    v.DeviceName,
			})
		}
	} else {
		if in.VersionList != nil {
			di, err := l.DiDB.FindByFilter(l.ctx, relationDB.DeviceFilter{
				ProductID: di.ProductID,
				Versions:  versionList,
			}, &def.PageInfo{Size: 100, Page: 1})
			if err != nil {
				return nil, err
			}
			for _, v := range di {
				otDB.Insert(l.ctx, &relationDB.DmOtaTaskDevices{
					FirmwareID:    in.FirmwareID,
					TaskUid:       taskUid,
					ProductID:     productFirmwareInfo.ProductID,
					Status:        101,
					TargetVersion: v.Version,
					DeviceName:    v.DeviceName,
				})
			}
		}
		if in.DeviceList != nil {
			for _, v := range deviceList {
				otDB.Insert(l.ctx, &relationDB.DmOtaTaskDevices{
					FirmwareID: in.FirmwareID,
					TaskUid:    taskUid,
					ProductID:  productFirmwareInfo.ProductID,
					Status:     101,
					DeviceName: v,
				})
			}
		}
	}

	return &dm.OtaTaskCreatResp{TaskID: taskID}, nil
}
