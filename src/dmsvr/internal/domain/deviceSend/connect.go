package deviceSend

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/repo/mysql"
	"github.com/go-things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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

func (l *ConnectLogic) Handle(msg *Elements) error {
	l.Infof("ConnectLogic|req=%+v", msg)
	ld, err := dm.GetClientIDInfo(msg.ClientID)
	_, _ = l.svcCtx.DeviceLog.Insert(&mysql.DeviceLog{
		ProductID:   ld.ProductID,
		Action:      msg.Action,
		Timestamp:   time.Unix(msg.Timestamp, 0), // 操作时间
		DeviceName:  ld.DeviceName,
		TranceID:    utils.TraceIdFromContext(l.ctx),
		Content:     msg.Payload,
		Topic:       msg.Topic,
		ResultType:  errors.Fmt(err).GetCode(),
		CreatedTime: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}
