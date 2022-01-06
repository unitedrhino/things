package logic

import (
	"context"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/exchange/types"
	"github.com/go-things/things/src/dmsvr/internal/repo/model/mysql"
	"github.com/go-things/things/src/dmsvr/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
	"time"
)

type ConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) LogicHandle {
	return LogicHandle(&ConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

func (l *ConnectLogic) Handle(msg *types.Elements) error {
	l.Infof("ConnectLogic|req=%+v", msg)
	ld, err := dm.GetClientIDInfo(msg.ClientID)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.DeviceLog.Insert(mysql.DeviceLog{
		ProductID:   ld.ProductID,
		Action:      msg.Action,
		Timestamp:   time.UnixMilli(msg.Timestamp), // 操作时间
		DeviceName:  ld.DeviceName,
		Payload:     msg.Payload,
		Topic:       msg.Topic,
		CreatedTime: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}
