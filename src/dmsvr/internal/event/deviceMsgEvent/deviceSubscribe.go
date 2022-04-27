package deviceMsgEvent

//设备的发布,连接及断连处理
import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceMsgHandle struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewDeviceMsgHandle(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceMsgHandle {
	return &DeviceMsgHandle{
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
	}
}

func (l *DeviceMsgHandle) Publish(msg *device.PublishMsg) error {
	l.Infof("DevReqLogic|req=%+v", utils.GetJson(msg))
	return NewPublishLogic(l.ctx, l.svcCtx).Handle(msg)
}

func (l *DeviceMsgHandle) Connected(msg *device.ConnectMsg) error {
	l.Infof("ConnectLogic|req=%+v", utils.GetJson(msg))
	//todo 这里需要查询下数据库,避免数据错误
	ld, err := device.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	err = l.svcCtx.DeviceLogRepo.Insert(l.ctx, &device.Log{
		ProductID:  ld.ProductID,
		Action:     msg.Action,
		Timestamp:  msg.Timestamp, // 操作时间
		DeviceName: ld.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
	})
	if err != nil {
		l.Errorf("%s|LogRepo|insert|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}

	return nil
}

func (l *DeviceMsgHandle) Disconnected(msg *device.ConnectMsg) error {
	l.Infof("DisconnectLogic|req=%+v", utils.GetJson(msg))
	ld, err := device.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	err = l.svcCtx.DeviceLogRepo.Insert(l.ctx, &device.Log{
		ProductID:  ld.ProductID,
		Action:     msg.Action,
		Timestamp:  msg.Timestamp, // 操作时间
		DeviceName: ld.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
	})
	if err != nil {
		l.Errorf("%s|LogRepo|insert|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	return nil
}
