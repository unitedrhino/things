package common

import (
	"context"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type InitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitLogic {
	return &InitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitLogic) Init() error {
	_, err := l.svcCtx.SchemaM.CommonSchemaInit(l.ctx, &dm.Empty{})
	return err
}
