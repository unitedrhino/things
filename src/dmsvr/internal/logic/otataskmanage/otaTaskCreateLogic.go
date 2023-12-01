package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskCreateLogic {
	return &OtaTaskCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 创建批量升级任务
func (l *OtaTaskCreateLogic) OtaTaskCreate(in *dm.OtaTaskCreateReq) (*dm.OtaTaskCreatResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaTaskCreatResp{}, nil
}
