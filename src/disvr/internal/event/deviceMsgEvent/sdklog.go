package deviceMsgEvent

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type SDKLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	topics []string
	dreq   deviceSend.DeviceReq
}

func NewSDKLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SDKLogLogic {
	return &SDKLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SDKLogLogic) initMsg(msg *deviceMsg.PublishMsg) error {
	var err error
	if err != nil {
		return err
	}
	err = utils.Unmarshal([]byte(msg.Payload), &l.dreq)
	if err != nil {
		return errors.Parameter.AddDetail("things topic is err:" + msg.Topic)
	}
	l.topics = strings.Split(msg.Topic, "/")
	return nil
}

func (l *SDKLogLogic) Handle(msg *deviceMsg.PublishMsg) (err error) {
	l.Infof("%s|req=%+v", utils.FuncName(), msg)
	err = l.initMsg(msg)
	if err != nil {
		return err
	}
	switch l.dreq.Type {
	case "get_log_level":
		l.GetLogLevel(msg)
	case "report_log_content":
		l.ReportLogContent(msg)
	}

	return nil
}

//获取设备上传的调试日志内容
func (l *SDKLogLogic) ReportLogContent(msg *deviceMsg.PublishMsg) error {
	ld, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	})
	if err != nil {
		l.Errorf("%s|Log|operate|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		return err
	}
	logContent := l.dreq.Params["content"]
	err = l.svcCtx.SDKLogRepo.Insert(l.ctx, &deviceMsg.SDKLog{
		ProductID:   ld.ProductID,
		LogLevel:    ld.LogLevel,
		Timestamp:   msg.Timestamp, // 操作时间
		DeviceName:  ld.DeviceName,
		Content:     logContent.(string),
		ClientToken: l.dreq.ClientToken,
	})
	if err != nil {
		l.Errorf("%s|LogRepo|insert|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		l.DeviceResp(msg, errors.Database, nil)
		return err
	}
	l.DeviceResp(msg, errors.OK, nil)
	return nil
}

//获取当前日志等级 0 未开启
func (l *SDKLogLogic) GetLogLevel(msg *deviceMsg.PublishMsg) {
	ld, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	})
	if err != nil {
		l.Errorf("%s|Log|operate|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		l.DeviceResp(msg, errors.Database, nil)
		return
	}
	l.DeviceResp(msg, errors.OK, map[string]any{"log_level": ld.LogLevel})
}

func (l *SDKLogLogic) DeviceResp(msg *deviceMsg.PublishMsg, err error, data map[string]any) {
	topic, payload := deviceSend.GenThingDeviceRespData(l.dreq.Method, l.dreq.ClientToken, l.topics, err, data)
	er := l.svcCtx.PubDev.PublishToDev(l.ctx, topic, payload)
	if er != nil {
		l.Errorf("DeviceResp|PublishToDev failure err:%v", er)
		return
	}
	l.Infof("ThingLogic|DeviceResp|topic:%v payload:%v", topic, string(payload))
}
