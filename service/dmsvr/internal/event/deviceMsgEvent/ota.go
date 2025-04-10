package deviceMsgEvent

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/application"
	"gitee.com/unitedrhino/things/share/userSubscribe"
	"github.com/spf13/cast"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgOta"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	//topics []string
	dreq msgOta.Req
	preq msgOta.Process
}

func NewOtaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaLogic {
	return &OtaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *OtaLogic) initMsg(msg *deviceMsg.PublishMsg) error {

	return nil
}
func (l *OtaLogic) Handle(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), msg)
	err = l.initMsg(msg)
	if err != nil {
		return nil, err
	}
	data, err := func() (data any, err error) {
		switch msg.Type {
		//case msgOta.TypeUpgrade: //固件升级消息上行 上报版本、拉取升级包
		//	return l.HandleUpgrade(msg)
		case msgOta.TypeProgress: //设备端上报升级进度
			return nil, l.HandleProgress(msg)
		default:
			return nil, errors.Parameter.AddDetail("things topic is err:" + msg.Topic)
		}
	}()
	respMsg = l.DeviceResp(msg, err, data)
	if l.dreq.NoAsk() { //如果不需要回复
		respMsg = nil
	}
	hub := &deviceLog.Hub{
		ProductID:   msg.ProductID,
		Action:      deviceLog.ActionTypeOta,
		Timestamp:   l.dreq.GetTimeStamp(), // 操作时间
		DeviceName:  msg.DeviceName,
		TraceID:     utils.TraceIdFromContext(l.ctx),
		RequestID:   l.dreq.MsgToken,
		Content:     string(msg.Payload),
		Topic:       msg.Topic,
		ResultCode:  errors.Fmt(err).GetCode(),
		RespPayload: respMsg.GetPayload(),
	}
	l.svcCtx.HubLogRepo.Insert(l.ctx, hub)
	l.svcCtx.UserSubscribe.Publish(l.ctx, userSubscribe.DevicePublish, hub.ToApp(), map[string]any{
		"productID":  msg.ProductID,
		"deviceName": msg.DeviceName,
	})
	return
}

//
//// 固件升级消息上行 上报版本、模块信息
//func (l *OtaLogic) HandleUpgrade(msg *deviceMsg.PublishMsg) (respData *msgOta.UpgradeData, err error) {
//	err = utils.Unmarshal([]byte(msg.Payload), &l.dreq)
//	if err != nil {
//		return nil, errors.Parameter.AddDetail("ota topic is err:" + msg.Topic)
//	}
//	err = l.dreq.VerifyReqParam()
//	if err != nil {
//		return nil, err
//	}
//	df, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
//		ProductID:    msg.ProductID,
//		DeviceNames:  []string{msg.DeviceName},
//		WithFirmware: true,
//		Statues:      []int64{msgOta.DeviceStatusInProgress, msgOta.DeviceStatusNotified, msgOta.DeviceStatusQueued},
//	})
//	if err != nil && !errors.Cmp(err, errors.NotFind) {
//		return nil, err
//	}
//	if df == nil {
//		jobs, err := relationDB.NewOtaJobRepo(l.ctx).FindByFilter(l.ctx, relationDB.OtaJobFilter{
//			ProductID:    msg.ProductID,
//			Statues:      []int64{msgOta.JobStatusInProgress},
//			UpgradeType:  msgOta.DynamicUpgrade, //静态升级需要先创建好设备,动态升级可以设备自己去获取
//			WithFirmware: true,
//			WithFiles:    true,
//		}, nil)
//		if err != nil {
//			return nil, err
//		}
//		for _, job := range jobs {
//			if utils.SliceIn(l.dreq.Params.Version, job.SrcVersions...) {
//				//如果在动态升级的版本内,则返回该升级包
//				df = &relationDB.DmOtaFirmwareDevice{
//					FirmwareID:  job.FirmwareID,
//					ProductID:   msg.ProductID,
//					DeviceName:  msg.DeviceName,
//					JobID:       job.ID,
//					SrcVersion:  l.dreq.Params.Version,
//					DestVersion: job.Firmware.Version,
//					Msg:      msgOta.DeviceStatusNotified,
//					Detail:      "设备主动拉取升级包",
//				}
//				err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Insert(l.ctx, df)
//				if err != nil {
//					return nil, err
//				}
//				df.Firmware = job.Firmware
//				df.Files = job.Files
//			}
//		}
//	}
//	return otamanagelogic.GenUpgradeParams(l.ctx, l.svcCtx, df.Firmware, df.Files)
//}

func (l *OtaLogic) HandleProgress(msg *deviceMsg.PublishMsg) (err error) {
	//设备端上报升级进度
	err = utils.Unmarshal([]byte(msg.Payload), &l.preq)
	if err != nil {
		return errors.Parameter.AddDetail("ota topic is err:" + msg.Topic)
	}
	err = l.preq.VerifyReqParam()
	if err != nil {
		return err
	}
	core := devices.Core{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, core)
	if err != nil {
		return err
	}
	df, err := relationDB.NewOtaFirmwareDeviceRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.OtaFirmwareDeviceFilter{
		ProductID:   msg.ProductID,
		DeviceNames: []string{msg.DeviceName},
		Statues:     []int64{msgOta.DeviceStatusInProgress, msgOta.DeviceStatusNotified},
	})
	if err != nil {
		return err
	}
	df.Step = l.preq.Params.Step
	df.Status = msgOta.DeviceStatusInProgress
	if l.preq.Params.Step < 0 {
		df.Status = msgOta.DeviceStatusFailure
		df.LastFailureTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		df.Detail = fmt.Sprintf("设备升级失败 上报错误:%d", l.preq.Params.Step)
	}

	err = relationDB.NewOtaFirmwareDeviceRepo(l.ctx).Update(l.ctx, df)
	if err != nil {
		return err
	}
	appMsg := application.OtaReport{
		Device:    core,
		Timestamp: time.Now().UnixMilli(),
		Status:    df.Status,
		Detail:    df.Detail,
		Step:      df.Step,
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
	return
}
func (l *OtaLogic) DeviceResp(msg *deviceMsg.PublishMsg, err error, data any) *deviceMsg.PublishMsg {
	if msg == nil {
		return nil
	}
	resp := &deviceMsg.CommonMsg{
		Method:   deviceMsg.GetRespMethod(l.dreq.Method),
		MsgToken: l.dreq.MsgToken,
		//Timestamp: time.Now().UnixMilli(),
		Data: data,
	}

	return &deviceMsg.PublishMsg{
		Handle:       msg.Handle,
		Type:         msg.Type,
		Payload:      resp.AddStatus(err, l.dreq.NeedRetMsg()).Bytes(),
		ProductID:    msg.ProductID,
		DeviceName:   msg.DeviceName,
		ProtocolCode: msg.ProtocolCode,
	}
}
