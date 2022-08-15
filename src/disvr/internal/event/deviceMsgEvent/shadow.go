package deviceMsgEvent

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ShadowLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShadowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShadowLogic {
	return &ShadowLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShadowLogic) Handle(msg *deviceMsg.PublishMsg) (err error) {
	l.Infof("%s|req=%+v", utils.FuncName(), msg)
	// todo
	return err
}
