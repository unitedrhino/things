package otataskmanagelogic

import (
	"context"

	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type OtaTaskReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOtaTaskReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OtaTaskReadLogic {
	return &OtaTaskReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 升级任务详情
func (l *OtaTaskReadLogic) OtaTaskRead(in *dm.OtaTaskReadReq) (*dm.OtaTaskReadResp, error) {
	otd, err := relationDB.NewOtaTaskRepo(l.ctx).FindOne(l.ctx, in.TaskID)
	if err != nil {
		return nil, err
	}
	return ToOtaTaskReadResp(otd), nil
}
