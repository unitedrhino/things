package deviceMsgEvent

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/service/deviceSend"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConnectedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	template *schema.Model
	topics   []string
	dreq     deviceSend.DeviceReq
}

func NewConnectedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectedLogic {
	return &ConnectedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
func (l *ConnectedLogic) Handle(msg *deviceMsg.ConnectMsg) error {
	l.Infof("%s req=%+v", utils.FuncName(), utils.Fmt(msg))
	ld, err := deviceAuth.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	err = l.svcCtx.HubLogRepo.Insert(l.ctx, &deviceMsg.HubLog{
		ProductID:  ld.ProductID,
		Action:     "connected",
		Timestamp:  msg.Timestamp, // 操作时间
		DeviceName: ld.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
	})
	if err != nil {
		l.Errorf("%s.LogRepo.insert productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	//更新对应设备的online状态
	di, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx, &dm.DeviceInfoReadReq{
		ProductID:  ld.ProductID,
		DeviceName: ld.DeviceName,
	})
	if err != nil {
		return err
	}
	l.svcCtx.DeviceM.DeviceInfoUpdate(l.ctx, &dm.DeviceInfo{
		ProductID:  di.ProductID,
		DeviceName: di.DeviceName,
		IsOnline:   def.True,
	})
	return nil
}
