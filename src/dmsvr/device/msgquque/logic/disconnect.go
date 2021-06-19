package logic

import (
	"context"
	"gitee.com/godLei6/things/src/dmsvr/device/msgquque/msvc"
	"gitee.com/godLei6/things/src/dmsvr/device/msgquque/types"
	"github.com/tal-tech/go-zero/core/logx"
)

type DisConnectLogic struct {
	ctx    context.Context
	svcCtx *msvc.ServiceContext
	logx.Logger
}

func NewDisConnectLogic(ctx context.Context, svcCtx *msvc.ServiceContext) LogicHandle {
	return LogicHandle(&DisConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

func (l *DisConnectLogic) Handle(msg *types.Elements) error {
	l.Infof("DisConnectLogic|req=%+v",msg)
	err :=  l.svcCtx.LogHandle(msg)
	if err != nil {
		return err
	}
	return nil
}