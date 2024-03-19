package devicemsglogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"time"

	"gitee.com/i-Things/share/domain/deviceMsg"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaPromptIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaPromptIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaPromptIndexLogic {
	return &OtaPromptIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ota
func (l *OtaPromptIndexLogic) OtaPromptIndex(in *dm.OtaPromptIndexReq) (*dm.OtaPromptIndexResp, error) {
	//var topic string
	////获取当前设备可用升级包
	//if in.GetId() > 0 && in.GetProductID() == "" {
	//	//主动重新升级
	//	taskDeviceInfo, err := otataskmanage.NewOtaTaskManageServer(l.svcCtx).OtaTaskDeviceRead(l.ctx, &dm.OtaTaskDeviceReadReq{
	//		ID: in.GetId(),
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//	if taskDeviceInfo.Status == 401 || taskDeviceInfo.Status == 501 {
	//		//升级中和升级成功的不能重新升级
	//		return nil, errors.OtaRetryStatusError
	//	}
	//	topic = fmt.Sprintf("$ota/down/upgrade/%s/%s", taskDeviceInfo.ProductID, taskDeviceInfo.DeviceName)
	//} else {
	//	topic = fmt.Sprintf("$ota/down/upgrade/%s/%s", in.ProductID, in.DeviceName)
	//}
	//
	////获取当前设备可用升级包
	//dmOtd := &dm.OtaTaskBatchReq{
	//	ID:         in.GetId(),
	//	ProductID:  in.GetProductID(),
	//	DeviceName: in.GetDeviceName(),
	//}
	//otd, err := otataskmanage.NewOtaTaskManageServer(l.svcCtx).OtaTaskDeviceEnableBatch(l.ctx, dmOtd)
	//if err != nil {
	//	//没找到可执行的升级任务
	//	return nil, err
	//}
	//firmwareInfo, err := relationDB.NewOtaFirmwareInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.OtaFirmwareInfoFilter{
	//	ID:        otd.FirmwareID,
	//	WithFiles: true,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//data := map[string]any{
	//	"id":         otd.ID,
	//	"version":    firmwareInfo.Version,
	//	"signMethod": firmwareInfo.SignMethod,
	//	"isDiff":     firmwareInfo.IsDiff,
	//}
	//files := make([]map[string]any, len(firmwareInfo.Files))
	//for k, v := range firmwareInfo.Files {
	//	url, _ := l.svcCtx.OssClient.PrivateBucket().SignedGetUrl(l.ctx, v.FilePath, 3600*24, common.OptionKv{})
	//	files[k] = map[string]any{
	//		"size":      v.Size,
	//		"signature": v.Signature,
	//		"name":      v.Name,
	//		"url":       url,
	//		//	"signMethod": v.SignMethod,
	//	}
	//}
	//data["files"] = files
	//dmReq := &dm.SendMsgReq{
	//	Topic: topic,
	//}
	//msgResp, err := server.NewDeviceInteractServer(l.svcCtx).SendMsg(l.ctx, l.DeviceResp(dmReq, errors.OK, data))
	////TODO 如何实时获取通知结果？？？
	//if err != nil {
	//	er := errors.Fmt(err)
	//	l.Errorf("%s.rpc.SendMsg req=%v err=%+v", utils.FuncName(), msgResp, er)
	//	return nil, err
	//}

	return &dm.OtaPromptIndexResp{}, nil
}
func (l *OtaPromptIndexLogic) DeviceResp(msg *dm.SendMsgReq, err error, data any) *dm.SendMsgReq {
	MsgToken := devices.GenMsgToken(l.ctx)
	resp := &deviceMsg.CommonMsg{
		Method:    "reportInfo",
		MsgToken:  MsgToken,
		Timestamp: time.Now().UnixMilli(),
		Data:      data,
	}

	return &dm.SendMsgReq{
		Topic:   msg.Topic,
		Payload: resp.AddStatus(err).Bytes(),
	}
}
