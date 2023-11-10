package deviceMsgEvent

import (
	"context"
	"strings"
	"time"

	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/oss/common"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgOta"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	topics []string
	dreq   msgOta.Req
	preq   msgOta.Process
}

func NewOtaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaLogic {
	return &OtaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *OtaLogic) initMsg(msg *deviceMsg.PublishMsg) error {
	l.topics = strings.Split(msg.Topic, "/")
	if len(l.topics) < 5 || l.topics[1] != "up" {
		return errors.Parameter.AddDetail("initMsg topic is err:" + msg.Topic)
	}
	return nil
}
func (l *OtaLogic) Handle(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), msg)
	err = l.initMsg(msg)
	if err != nil {
		return nil, err
	}
	respMsg, err = func() (respMsg *deviceMsg.PublishMsg, err error) {
		switch l.topics[2] {
		case msgOta.TypeReport: //固件升级消息上行 上报版本、模块信息
			return l.HandleReport(msg)
		case msgOta.TypeProgress: //设备端上报升级进度
			return l.HandleResp(msg)
		default:
			return nil, errors.Parameter.AddDetail("things topic is err:" + msg.Topic)
		}
	}()
	l.svcCtx.HubLogRepo.Insert(l.ctx, &msgHubLog.HubLog{
		ProductID:  msg.ProductID,
		Action:     "otaLog",
		Timestamp:  l.dreq.GetTimeStamp(), // 操作时间
		DeviceName: msg.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		RequestID:  l.dreq.MsgToken,
		Content:    string(msg.Payload),
		Topic:      msg.Topic,
		ResultType: errors.Fmt(err).GetCode(),
	})
	return
}
func (l *OtaLogic) HandleReport(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	//固件升级消息上行 上报版本、模块信息
	err = utils.Unmarshal([]byte(msg.Payload), &l.dreq)
	if err != nil {
		return nil, errors.Parameter.AddDetail("ota topic is err:" + msg.Topic)
	}
	err = l.dreq.VerifyReqParam()
	if err != nil {
		return nil, err
	}
	//获取当前设备可用升级包
	di := &dm.OtaTaskBatchReq{
		ID:         l.dreq.Params.ID,
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
		Module:     "", //default
		Version:    l.dreq.GetVersion(),
	}
	otd, err := l.svcCtx.OtaTaskM.OtaTaskDeviceEnableBatch(l.ctx, di)
	if err != nil {
		//没找到可执行的升级任务
		return nil, err
	}
	firmwareInfo, err := l.svcCtx.FirmwareM.FirmwareInfoRead(l.ctx, &dm.FirmwareInfoReadReq{
		FirmwareID: otd.FirmwareID,
	})
	if err != nil {
		return nil, err
	}
	data := map[string]any{
		"id":         otd.ID,
		"version":    firmwareInfo.Version,
		"signMethod": firmwareInfo.SignMethod,
		"isDiff":     firmwareInfo.IsDiff,
	}
	var files = make([]map[string]any, len(firmwareInfo.Files))
	for k, v := range firmwareInfo.Files {
		url, _ := l.svcCtx.OssClient.PrivateBucket().SignedGetUrl(l.ctx, v.FilePath, 3600*24, common.OptionKv{})
		files[k] = map[string]any{
			"size":      v.Size,
			"name":      v.Name,
			"url":       url,
			"signature": v.Signature,
		}
	}
	data["files"] = files
	return l.DeviceResp(msg, errors.OK, data), nil
}
func (l *OtaLogic) HandleResp(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	//设备端上报升级进度
	err = utils.Unmarshal([]byte(msg.Payload), &l.preq)
	if err != nil {
		return nil, errors.Parameter.AddDetail("ota topic is err:" + msg.Topic)
	}
	err = l.preq.VerifyReqParam()
	if err != nil {
		return nil, err
	}
	_, err = l.svcCtx.OtaTaskM.OtaTaskDeviceProcess(l.ctx, &dm.OtaTaskDeviceProcessReq{
		ID:     l.preq.Params.ID,
		Module: l.preq.Params.Module,
		Step:   l.preq.Params.Step,
		Desc:   l.preq.Params.Desc,
	})
	if err != nil {
		return nil, errors.Parameter.AddDetail("ota process rpc is err:" + msg.Topic)
	}
	return l.DeviceResp(msg, errors.OK, nil), nil
}
func (l *OtaLogic) DeviceResp(msg *deviceMsg.PublishMsg, err error, data any) *deviceMsg.PublishMsg {
	resp := &deviceMsg.CommonMsg{
		Method:    deviceMsg.GetRespMethod(l.dreq.Method),
		MsgToken:  l.dreq.MsgToken,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}
	if l.topics[2] == "report" {
		l.topics[2] = "upgrade" //下行 升级包详情
	}
	return &deviceMsg.PublishMsg{
		Handle:     msg.Handle,
		Type:       l.topics[2],
		Payload:    resp.AddStatus(err).Bytes(),
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}
}
