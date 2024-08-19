package deviceMsgEvent

import (
	"context"
	"gitee.com/i-Things/share/domain/deviceAuth"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceStatus"
	devicemanagelogic "github.com/i-Things/things/service/dmsvr/internal/logic/devicemanage"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type DisconnectedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	template *schema.Model
	topics   []string
	dreq     msgThing.Req
}

func NewDisconnectedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DisconnectedLogic {
	return &DisconnectedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *DisconnectedLogic) Handle(msg *deviceStatus.ConnectMsg) error {
	l.Infof("%s req=%+v", utils.FuncName(), utils.Fmt(msg))
	dev := msg.Device
	ld, err := deviceAuth.GetClientIDInfo(msg.ClientID)
	if err != nil && dev.DeviceName == "" {
		return err
	}
	if ld != nil && ld.IsNeedRegister {
		return nil
	}
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, dev)
	if err != nil {
		return err
	}
	if di.FirstLogin == 0 {
		return nil
	}
	err = devicemanagelogic.HandleOnlineFix(l.ctx, l.svcCtx, msg)
	if err != nil {
		l.Error(err)
	}
	//err = l.svcCtx.DeviceStatus.AddDevice(l.ctx, msg)
	return err
	////更新对应设备的online状态
	//di, err := relationDB.NewDeviceInfoRepo(l.ctx).FindOneByFilter(l.ctx, relationDB.DeviceFilter{
	//	ProductID:  ld.ProductID,
	//	DeviceName: ld.DeviceName,
	//})
	//if err != nil {
	//	l.Errorf("%s.DeviceStatusDisConnected productID:%v deviceName:%v err:%v",
	//		utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	//	return err
	//} else {
	//	if di.IsOnline != def.False {
	//		di.IsOnline = def.False
	//		err = relationDB.NewDeviceInfoRepo(l.ctx).Update(l.ctx, di)
	//		if err != nil {
	//			l.Errorf("%s.DeviceInfoUpdate productID:%v deviceName:%v err:%v",
	//				utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	//		}
	//	}
	//}
	//
	//err = l.svcCtx.StatusRepo.Insert(l.ctx, &deviceLog.Status{
	//	ProductID:  ld.ProductID,
	//	Status:     def.DisConnectedStatus,
	//	Timestamp:  msg.Timestamp, // 操作时间
	//	DeviceName: ld.DeviceName,
	//})
	//if err != nil {
	//	l.Errorf("%s.LogRepo.insert productID:%v deviceName:%v err:%v",
	//		utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	//}
	//err = l.svcCtx.PubApp.DeviceStatusDisConnected(l.ctx, application.ConnectMsg{
	//	Device: devices.Core{
	//		ProductID:  ld.ProductID,
	//		DeviceName: ld.DeviceName,
	//	},
	//	Timestamp: msg.Timestamp.UnixMilli(),
	//})
	//if err != nil {
	//	l.Errorf("%s.DeviceStatusDisConnected productID:%v deviceName:%v err:%v",
	//		utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	//}
	//
	//return nil
}
