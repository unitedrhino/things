package logic

import (
"context"
"gitee.com/godLei6/things/src/dmsvr/device/msgquque/msvc"
	"gitee.com/godLei6/things/src/dmsvr/device/msgquque/types"
	"github.com/tal-tech/go-zero/core/logx"
)

type PublishLogic struct {
	ctx    context.Context
	svcCtx *msvc.ServiceContext
	logx.Logger
}

func NewPublishLogic(ctx context.Context, svcCtx *msvc.ServiceContext) LogicHandle {
	return LogicHandle(&PublishLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	})
}

func (l *PublishLogic) Handle(msg *types.Elements) error {
	l.Infof("PublishLogic|req=%+v",msg)
	return nil
}