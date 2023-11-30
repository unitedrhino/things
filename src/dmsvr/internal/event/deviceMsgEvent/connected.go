package deviceMsgEvent

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/domain/application"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgHubLog"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceStatus"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type ConnectedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	template *schema.Model
	topics   []string
	dreq     msgThing.Req
	DiDB     *relationDB.DeviceInfoRepo
	di       *relationDB.DmDeviceInfo
}

func NewConnectedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConnectedLogic {
	return &ConnectedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
		Logger: logx.WithContext(ctx),
	}
}

func (l *ConnectedLogic) UpdateLoginTime() {
	now := sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	if l.di.FirstLogin.Valid == false {
		l.di.FirstLogin = now
	}
	l.di.LastLogin = now
	l.di.IsOnline = def.True
	l.DiDB.Update(l.ctx, l.di)
}

func (l *ConnectedLogic) Handle(msg *deviceStatus.ConnectMsg) error {
	l.Infof("%s req=%+v", utils.FuncName(), utils.Fmt(msg))
	ld, err := deviceAuth.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	l.di, err = l.DiDB.FindOneByFilter(l.ctx, relationDB.DeviceFilter{ProductID: ld.ProductID, DeviceNames: []string{ld.DeviceName}})
	if err != nil {
		return err
	}
	l.UpdateLoginTime()
	err = l.svcCtx.HubLogRepo.Insert(l.ctx, &msgHubLog.HubLog{
		ProductID:  ld.ProductID,
		Action:     deviceStatus.ConnectStatus,
		Timestamp:  msg.Timestamp, // 操作时间
		DeviceName: ld.DeviceName,
		TranceID:   utils.TraceIdFromContext(l.ctx),
		ResultType: errors.Fmt(err).GetCode(),
	})
	if err != nil {
		l.Errorf("%s.HubLogRepo.insert productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	err = l.svcCtx.PubApp.DeviceStatusConnected(l.ctx, application.ConnectMsg{
		Device: devices.Core{
			ProductID:  ld.ProductID,
			DeviceName: ld.DeviceName,
		},
		Timestamp: msg.Timestamp.UnixMilli(),
	})
	if err != nil {
		l.Errorf("%s.PubApp.DeviceStatusConnected productID:%v deviceName:%v err:%v",
			utils.FuncName(), ld.ProductID, ld.DeviceName, err)
	}
	return nil
}
