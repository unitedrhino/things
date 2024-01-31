package otataskmanagelogic

import (
	"context"

	"gitee.com/i-Things/core/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
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
	var taskDB = relationDB.NewOtaTaskRepo(l.ctx)
	//ota批量取消
	otd, err := taskDB.FindOne(l.ctx, in.TaskID)
	if err != nil {
		return nil, err
	}
	if otd.Status > 2 {
		return nil, errors.OtaCancleStatusError
	}
	otd.Status = 4
	err = taskDB.Update(l.ctx, otd)
	if err == nil {
		err = relationDB.NewOtaTaskDevicesRepo(l.ctx).CancelByTaskUid(l.ctx, otd.TaskUid)
	}
	return &dm.OtaCommonResp{}, err
}
