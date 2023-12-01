package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskCancleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskCancleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskCancleLogic {
	return &OtaTaskCancleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量取消升级任务
func (l *OtaTaskCancleLogic) OtaTaskCancle(in *dm.OtaTaskCancleReq) (*dm.OtaCommonResp, error) {
	// todo: add your logic here and delete this line

	return &dm.OtaCommonResp{}, nil
}
