package deviceMsgEvent

import (
	"context"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/devices"
	"gitee.com/i-Things/core/shared/domain/application"
	"gitee.com/i-Things/core/shared/domain/deviceAuth"
	"gitee.com/i-Things/core/shared/domain/schema"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceStatus"
	server "github.com/i-Things/things/src/dmsvr/internal/server/devicemanage"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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
	ld, err := deviceAuth.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	err = l.svcCtx.HubLogRepo.Insert(l.ctx, &msgHubLog.HubLog{
		ProductID:  ld.ProductID,
		Action:     deviceStatus.DisConnectStatus,
		Timestamp:  msg.Timestamp, // 操作时间
		DeviceName: ld.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
	})
	if err != nil {
		l.Errorf("%s.LogRepo.insert productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	err = l.svcCtx.PubApp.DeviceStatusDisConnected(l.ctx, application.ConnectMsg{
		Device: devices.Core{
			ProductID:  ld.ProductID,
			DeviceName: ld.DeviceName,
		},
		Timestamp: msg.Timestamp.UnixMilli(),
	})
	if err != nil {
		l.Errorf("%s.DeviceStatusDisConnected productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	//更新对应设备的online状态
	_, err = server.NewDeviceManageServer(l.svcCtx).DeviceInfoUpdate(l.ctx, &dm.DeviceInfo{
		ProductID:  ld.ProductID,
		DeviceName: ld.DeviceName,
		IsOnline:   def.False,
	})
	if err != nil {
		l.Errorf("%s.DeviceInfoUpdate productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	return nil
}
