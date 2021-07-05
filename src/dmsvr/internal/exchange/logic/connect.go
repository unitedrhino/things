package logic

import (
	"context"
	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"gitee.com/godLei6/things/src/dmsvr/model"
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
	_, err = l.svcCtx.DeviceLog.Insert(model.DeviceLog{
		ProductID:   ld.ProductID,
		Action:      msg.Action,
		Timestamp:   time.Unix(msg.Timestamp, 0), // 操作时间
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
