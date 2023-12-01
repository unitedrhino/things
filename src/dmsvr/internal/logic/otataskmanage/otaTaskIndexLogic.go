package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskIndexLogic {
	return &OtaTaskIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaTaskIndexLogic) OtaTaskIndex(in *dm.OtaTaskIndexReq) (*dm.OtaTaskIndexResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaTaskIndexResp{}, nil
}
