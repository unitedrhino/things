package logic

import (
	"context"
	"gitee.com/godLei6/things/src/dmsvr/internal/msgquque/msvc"
	"gitee.com/godLei6/things/src/dmsvr/internal/msgquque/types"
	"github.com/tal-tech/go-zero/core/logx"
)

type ConnectLogic struct {
	ctx    context.Context
	svcCtx *msvc.ServiceContext
	logx.Logger
}

func NewConnectLogic(ctx context.Context, svcCtx *msvc.ServiceContext) LogicHandle {
	return LogicHandle(&ConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

func (l *ConnectLogic) Handle(msg *types.Elements) error {
	l.Infof("ConnectLogic|req=%+v",msg)

	return nil
}