package deviceMsgEvent

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceData"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/dmsvr/internal/domain/thing"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConnectedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	template *thing.Template
	topics   []string
	dreq     deviceSend.DeviceReq
	dd       deviceData.DeviceDataRepo
}

func NewConnectedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectedLogic {
	return &ConnectedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ConnectedLogic) Handle(msg *device.ConnectMsg) error {
	l.Infof("%s|req=%+v", utils.FuncName(), utils.GetJson(msg))
	ld, err := device.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	err = l.svcCtx.HubLogRepo.Insert(l.ctx, &device.HubLog{
		ProductID:  ld.ProductID,
		Action:     "connected",
		Timestamp:  msg.Timestamp, // 操作时间
		DeviceName: ld.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
	})
	if err != nil {
		l.Errorf("%s|LogRepo|insert|productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	//更新对应设备的online状态
	di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, ld.ProductID, ld.DeviceName)
	if err != nil {
		if err == mysql.ErrNotFound {
			return errors.NotFind.AddDetailf("Disconnect|not find device|productid=%s|deviceName=%s",
				ld.ProductID, ld.DeviceName)
		}
		return errors.Database.AddDetail(err.Error())
	}
	di.IsOnline = 1 //在线
	l.svcCtx.DeviceInfo.Update(l.ctx, di)
	return nil
}
