package logic

import (
	"context"
	"gitee.com/godLei6/things/src/dmsvr/internal/exchange/types"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"
	"github.com/tal-tech/go-zero/core/logx"
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

func (l *DisConnectLogic) Handle(msg *types.Elements) error {
	l.Infof("DisConnectLogic|req=%+v", msg)

	return nil
}
