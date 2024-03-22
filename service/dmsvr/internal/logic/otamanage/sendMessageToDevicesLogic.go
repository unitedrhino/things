package otamanagelogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/domain/deviceMsg"
	"gitee.com/i-Things/share/domain/deviceMsg/msgOta"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/oss/common"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type SendMessageToDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB  *relationDB.DeviceInfoRepo
	GdDB  *relationDB.GroupDeviceRepo
	OtDB  *relationDB.OtaFirmwareDeviceRepo
	OjDB  *relationDB.OtaJobRepo
	OfDB  *relationDB.OtaFirmwareInfoRepo
	OffDB *relationDB.OtaFirmwareFileRepo
}

func NewSendMessageToDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageToDevicesLogic {
	return &SendMessageToDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
		OtDB:   relationDB.NewOtaFirmwareDeviceRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareInfoRepo(ctx),
		OffDB:  relationDB.NewOtaFirmwareFileRepo(ctx),
	}
}

func (l *SendMessageToDevicesLogic) PushMessageToDevices(jobInfo *relationDB.DmOtaFirmwareJob) error {
	firmware := jobInfo.Firmware
	if jobInfo.IsNeedPush != def.True { //只有需要推送的才推送
		return nil
	}
	deviceList, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
		FirmwareID: jobInfo.FirmwareID,
		JobID:      jobInfo.ID,
		ProductID:  firmware.ProductID,
		IsOnline:   def.True,                           //只有在线的设备才推送
		Statues:    []int64{msgOta.DeviceStatusQueued}, //只处理待推送的设备
	}, &def.PageInfo{
		Page: 1,
		Size: jobInfo.MaximumPerMinute/(60/5) + 1, //任务5秒钟推送一次
	})

	firmwareFiles, err := l.OffDB.FindByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{FirmwareID: jobInfo.FirmwareID}, nil)
	if err != nil {
		return err
	}
	data, err := GenUpgradeParams(l.ctx, l.svcCtx, firmware, firmwareFiles)
	if err != nil {
		return err
	}
	MsgToken := devices.GenMsgToken(l.ctx)
	upgradeMsg := deviceMsg.CommonMsg{
		MsgToken:  MsgToken,
		Method:    msgOta.TypeUpgrade,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}
	payload, _ := json.Marshal(upgradeMsg)
	pi, err := l.svcCtx.ProductCache.GetData(l.ctx, firmware.ProductID)
	if err != nil {
		return err
	}
	for _, device := range deviceList {
		reqMsg := deviceMsg.PublishMsg{
			Handle:       devices.Ota,
			Type:         msgOta.TypeUpgrade,
			Payload:      payload,
			Timestamp:    time.Now().UnixMilli(),
			ProductID:    device.ProductID,
			DeviceName:   device.DeviceName,
			ProtocolCode: pi.ProtocolCode,
		}
		err := l.svcCtx.PubDev.PublishToDev(l.ctx, &reqMsg)
		if err != nil {
			return err
		}
		device.Status = msgOta.DeviceStatusNotified
		device.Detail = "主动推送"
		err = relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, device)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenUpgradeParams(ctx context.Context, svcCtx *svc.ServiceContext, firmware *relationDB.DmOtaFirmwareInfo, files []*relationDB.DmOtaFirmwareFile) (*msgOta.UpgradeParams, error) {
	if len(files) == 0 {
		return nil, errors.System.AddDetail("升级包下没有文件")
	}
	if len(files) == 1 { //单文件模式
		url, err := svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, files[0].FilePath, 60*20, common.OptionKv{})
		if err != nil {
			return nil, err
		}
		data := msgOta.UpgradeParams{
			Version:    firmware.Version,
			IsDiff:     firmware.IsDiff,
			SignMethod: firmware.SignMethod,
			Extra:      firmware.Extra,
			File: &msgOta.File{
				Size:      files[0].Size,
				Name:      files[0].Name,
				FileUrl:   url,
				FileMd5:   files[0].FileMd5,
				Signature: files[0].Signature,
			},
		}
		return &data, nil
	}
	var data = msgOta.UpgradeParams{
		Version:    firmware.Version,
		IsDiff:     firmware.IsDiff,
		SignMethod: firmware.SignMethod,
		Extra:      firmware.Extra,
	}
	for _, f := range files {
		url, err := svcCtx.OssClient.PrivateBucket().SignedGetUrl(ctx, f.FilePath, 60*20, common.OptionKv{})
		if err != nil {
			return nil, err
		}
		data.Files = append(data.Files, &msgOta.File{
			Size:      f.Size,
			Name:      f.Name,
			FileUrl:   url,
			FileMd5:   f.FileMd5,
			Signature: f.Signature,
		})
	}
	return &data, nil
}
