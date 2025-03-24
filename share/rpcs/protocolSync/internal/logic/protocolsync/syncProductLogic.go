package protocolsynclogic

import (
	"context"

	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/internal/svc"
	"gitee.com/unitedrhino/things/share/rpcs/protocolSync/pb/protocolSync"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncProductLogic {
	return &SyncProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncProductLogic) SyncProduct(in *protocolSync.Empty) (*protocolSync.Empty, error) {
	// todo: add your logic here and delete this line

	return &protocolSync.Empty{}, nil
}
