package deviceMsgEvent

import (
	"context"
	"strings"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
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
		return errors.Parameter.AddDetail("sdklog topic is err:" + msg.Topic)
	}
	l.topics = strings.Split(msg.Topic, "/")
	return nil
}

func (l *SDKLogLogic) Handle(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), msg)
	err = l.initMsg(msg)
	if err != nil {
		return nil, err
	}
	switch l.dreq.Type {
	case "get_log_level":
		respMsg, err = l.GetLogLevel(msg)
	case "report_log_content":
		respMsg, err = l.ReportLogContent(msg)
	default:
		return nil, errors.Parameter.AddDetail("sdk log topic is err:" + msg.Topic)
	}
	return respMsg, err
}

//获取设备上传的调试日志内容
func (l *SDKLogLogic) ReportLogContent(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	ld, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	})
	if err != nil {
		l.Errorf("%s.Log.operate productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		return nil, err
	}
	logContent, ok := l.dreq.Params["content"]
	if !ok {
		return nil, errors.SDKLOG_MissParam
	}
	logLevel, ok := l.dreq.Params["log_level"]
	if !ok {
		logLevel = def.LogDebug
	} else {
		if logLevel.(int64) < int64(def.LogError) || logLevel.(int64) > def.LogDebug {
			return nil, errors.SDKLOG_ErrorLevel
		}
	}
	err = l.svcCtx.SDKLogRepo.Insert(l.ctx, &deviceMsg.SDKLog{
		ProductID:   ld.ProductID,
		LogLevel:    logLevel.(int64),
		Timestamp:   msg.Timestamp, // 操作时间
		DeviceName:  ld.DeviceName,
		Content:     logContent.(string),
		ClientToken: l.dreq.ClientToken,
	})
	if err != nil {
		l.Errorf("%s.LogRepo.insert.productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)

		return l.DeviceResp(msg, errors.Database, nil), err
	}

	return l.DeviceResp(msg, errors.OK, nil), nil
}

//获取当前日志等级 1 未开启
func (l *SDKLogLogic) GetLogLevel(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	ld, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	})
	if err != nil {
		l.Errorf("%s.Log.operate.productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
		return l.DeviceResp(msg, errors.Database, nil), err
	}
	return l.DeviceResp(msg, errors.OK, map[string]any{"log_level": ld.LogLevel}), nil
}

func (l *SDKLogLogic) DeviceResp(msg *deviceMsg.PublishMsg, err error, data map[string]any) *deviceMsg.PublishMsg {
	topic, payload := deviceSend.GenThingDeviceRespData(l.dreq.Method, l.dreq.ClientToken, l.topics, err, data)
	return &deviceMsg.PublishMsg{
		Topic:   topic,
		Payload: payload,
	}
}
