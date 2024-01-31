package otajobmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	server "github.com/i-Things/things/service/dmsvr/internal/server/deviceinteract"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type SendMessageToDevicesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB  *relationDB.DeviceInfoRepo
	GdDB  *relationDB.GroupDeviceRepo
	OtDB  *relationDB.OtaUpgradeTaskRepo
	OjDB  *relationDB.OtaJobRepo
	OfDB  *relationDB.OtaFirmwareRepo
	OffDB *relationDB.OtaFirmwareFileRepo
}

func NewSendMessageToDevicesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageToDevicesLogic {
	return &SendMessageToDevicesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		OjDB:   relationDB.NewOtaJobRepo(ctx),
		OtDB:   relationDB.NewOtaUpgradeTaskRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		OfDB:   relationDB.NewOtaFirmwareRepo(ctx),
		OffDB:  relationDB.NewOtaFirmwareFileRepo(ctx),
	}
}

func (l *SendMessageToDevicesLogic) PushMessageToDevices(deviceList []*relationDB.DmDeviceInfo, firmware *relationDB.DmOtaFirmware) error {
	firmwareFiles, err := l.OffDB.FindByFilter(l.ctx, relationDB.OtaFirmwareFileFilter{FirmwareID: firmware.ID}, nil)
	if err != nil {
		return err
	}
	var files []msgOta.File
	_ = copier.Copy(&files, &firmwareFiles)
	//签名值
	for k, file := range firmwareFiles {
		files[k].FileMd5 = file.Signature
	}

	MsgToken := devices.GenMsgToken(l.ctx)
	upgradeMsg := deviceMsg.CommonMsg{
		MsgToken:  MsgToken,
		Method:    msgOta.TypeUpgrade,
		Timestamp: time.Now().UnixMilli(),
		Data: msgOta.UpgradeParams{
			Version:          firmware.Version,
			IsDiff:           firmware.IsDiff,
			SignMethod:       firmware.SignMethod,
			DownloadProtocol: "https",
			Module:           firmware.Module,
			Files:            files,
		},
	}
	for _, device := range deviceList {
		topic := fmt.Sprintf("$ota/down/upgrade/%s/%s", firmware.ProductID, device.DeviceName)
		logx.Infof("topic:%+v", topic)
		sendMsg := dm.SendMsgReq{
			Topic:   topic,
			Payload: upgradeMsg.AddStatus(errors.OK).Bytes(),
		}
		logx.Infof("Payload: %q\n", sendMsg.Payload)
		logx.Infof("sendMsg:%+v", &sendMsg)
		_, err := server.NewDeviceInteractServer(l.svcCtx).SendMsg(l.ctx, &sendMsg)
		if err != nil {
			l.Errorf("错误是", err)
			logx.Infof("消息发送失败")
			return err
		}
	}
	return nil
}
