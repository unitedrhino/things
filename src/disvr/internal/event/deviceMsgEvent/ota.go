package deviceMsgEvent

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type OtaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaLogic {
	return &OtaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaLogic) Handle(msg *deviceMsg.PublishMsg) (err error) {
	l.Infof("%s req=%+v", utils.FuncName(), msg)
	// todo
	return err
}
