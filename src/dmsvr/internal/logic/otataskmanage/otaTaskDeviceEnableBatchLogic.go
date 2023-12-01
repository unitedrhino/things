package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskDeviceEnableBatchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskDeviceEnableBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskDeviceEnableBatchLogic {
	return &OtaTaskDeviceEnableBatchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取当前可执行批次信息
func (l *OtaTaskDeviceEnableBatchLogic) OtaTaskDeviceEnableBatch(in *dm.OtaTaskBatchReq) (*dm.OtaTaskBatchResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaTaskBatchResp{}, nil
}
