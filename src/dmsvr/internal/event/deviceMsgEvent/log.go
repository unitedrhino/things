package deviceMsgEvent

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type LogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	topics []string
	dreq   deviceSend.DeviceReq
}

func NewLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogLogic {
	return &LogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogLogic) initMsg(msg *device.PublishMsg) error {
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

func (l *LogLogic) Handle(msg *device.PublishMsg) (err error) {
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

func (l *LogLogic) ReportLogContent(msg *device.PublishMsg) error {
	ld, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(msg.ProductID, msg.DeviceName)
	if err != nil {
		l.Errorf("%s|Log|operate|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		return err
	}
	if err := l.svcCtx.DebugLogRepo.InitProduct(
		l.ctx, ld.ProductID); err != nil {
		l.Errorf("%s|DeviceLogRepo|InitProduct| failure,err:%v", utils.FuncName(), err)
		return errors.Database.AddDetail(err)
	}
	l.svcCtx.DebugLogRepo.InitDevice(l.ctx, ld.ProductID, ld.DeviceName)
	logContent := l.dreq.Params["content"]

	err = l.svcCtx.DebugLogRepo.Insert(l.ctx, &device.DebugLog{
		ProductID:  ld.ProductID,
		LogLevel:   ld.LogLevel,
		Timestamp:  msg.Timestamp, // 操作时间
		DeviceName: ld.DeviceName,
		Content:    logContent.(string),
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
		RequestID:  l.dreq.ClientToken,
	})
	if err != nil {
		l.Errorf("%s|LogRepo|insert|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		l.DeviceResp(msg, errors.Database, nil)
	}
	l.DeviceResp(msg, errors.OK, map[string]interface{}{"RequestId": l.dreq.ClientToken})
	return nil
}

func (l *LogLogic) GetLogLevel(msg *device.PublishMsg) {
	ld, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(msg.ProductID, msg.DeviceName)
	if err != nil {
		l.Errorf("%s|Log|operate|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		l.DeviceResp(msg, errors.Database, nil)
	}
	l.DeviceResp(msg, errors.OK, map[string]interface{}{"log_level": ld.LogLevel})
}

func (l *LogLogic) DeviceResp(msg *device.PublishMsg, err error, data map[string]interface{}) {
	topic, payload := deviceSend.GenThingDeviceRespData(l.dreq.Method, l.dreq.ClientToken, l.topics, err, data)
	er := l.svcCtx.InnerLink.PublishToDev(l.ctx, topic, payload)
	if er != nil {
		l.Errorf("DeviceResp|PublishToDev failure err:%v", er)
		return
	}
	l.Infof("ThingLogic|DeviceResp|topic:%v payload:%v", topic, string(payload))
}
