package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskUpdateLogic {
	return &OtaTaskUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OtaTaskUpdateLogic) OtaTaskUpdate(in *dm.OtaTaskInfo) (*dm.OtaCommonResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaCommonResp{}, nil
}
