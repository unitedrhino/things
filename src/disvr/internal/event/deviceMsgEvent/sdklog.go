package deviceMsgEvent

import (
	"context"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgSdkLog"
	"time"

	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type SDKLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	dreq msgSdkLog.Req
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
		return errors.Parameter.AddDetailf("sdkLog Unmarshal err:%v", err)
	}
	return nil
}

func (l *SDKLogLogic) Handle(msg *deviceMsg.PublishMsg) (respMsg *deviceMsg.PublishMsg, err error) {
	l.Infof("%s req=%+v", utils.FuncName(), msg)
	err = l.initMsg(msg)
	if err != nil {
		return nil, err
	}
	switch msg.Types[0] {
	case msgSdkLog.TypeOperation:
		respMsg, err = l.GetLogLevel(msg)
	case msgSdkLog.TypeReport:
		respMsg, err = l.ReportLogContent(msg)
	default:
		return nil, errors.Parameter.AddDetailf("sdk log types is err:%v", msg.Types)
	}
	l.svcCtx.HubLogRepo.Insert(l.ctx, &msgHubLog.HubLog{
		ProductID:  msg.ProductID,
		Action:     "sdkLog",
		Timestamp:  time.Now(), // 记录当前时间
		DeviceName: msg.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		RequestID:  l.dreq.ClientToken,
		Content:    string(msg.Payload),
		Topic:      msg.Handle, //todo 等待实现
		ResultType: errors.Fmt(err).GetCode(),
	})
	return respMsg, err
}

// 获取设备上传的调试日志内容
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
	err = l.dreq.VerifyReqParam()
	if err != nil {
		return nil, err
	}
	for _, logObj := range l.dreq.Params {
		err = l.svcCtx.SDKLogRepo.Insert(l.ctx, &msgSdkLog.SDKLog{
			ProductID:  ld.ProductID,
			LogLevel:   logObj.LogLevel,
			Timestamp:  l.dreq.GetTimeStamp(logObj.Timestamp), // 操作时间
			DeviceName: ld.DeviceName,
			Content:    logObj.Content,
		})
		if err != nil {
			l.Errorf("%s.LogRepo.insert.productID:%v deviceName:%v err:%v",
				utils.FuncName(), ld.ProductID, ld.DeviceName, err)

			return l.DeviceResp(msg, errors.Database, nil), err
		}
	}

	return l.DeviceResp(msg, errors.OK, nil), nil
}

// 获取当前日志等级
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
	return l.DeviceResp(msg, errors.OK, map[string]any{"logLevel": ld.LogLevel}), nil
}

func (l *SDKLogLogic) DeviceResp(msg *deviceMsg.PublishMsg, err error, data any) *deviceMsg.PublishMsg {
	resp := &deviceMsg.CommonMsg{
		Method:      deviceMsg.GetRespMethod(l.dreq.Method),
		ClientToken: l.dreq.ClientToken,
		Timestamp:   time.Now().UnixMilli(),
		Data:        data,
	}
	return &deviceMsg.PublishMsg{
		Handle:     msg.Handle,
		Types:      msg.Types,
		Payload:    resp.AddStatus(err).Bytes(),
		ProductID:  msg.ProductID,
		DeviceName: msg.DeviceName,
	}
}
