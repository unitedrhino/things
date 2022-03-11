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

type DisConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDisConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) LogicHandle {
	return LogicHandle(&DisConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

func (l *DisConnectLogic) Handle(msg *Elements) error {
	l.Infof("DisConnectLogic|req=%+v", msg)
	ld, err := dm.GetClientIDInfo(msg.ClientID)
	_, _ = l.svcCtx.DeviceLog.Insert(&mysql.DeviceLog{
		ProductID:   ld.ProductID,
		Action:      msg.Action,
		Timestamp:   time.Unix(msg.Timestamp, 0), // 操作时间
		DeviceName:  ld.DeviceName,
		TranceID:    utils.TraceIdFromContext(l.ctx),
		Content:     msg.Payload,
		ResultType:  errors.Fmt(err).GetCode(),
		Topic:       msg.Topic,
		CreatedTime: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}
