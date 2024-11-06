package devicemsglogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"time"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type AbnormalLogCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAbnormalLogCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AbnormalLogCreateLogic {
	return &AbnormalLogCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AbnormalLogCreateLogic) AbnormalLogCreate(in *dm.AbnormalLogInfo) (*dm.Empty, error) {
	if in.Timestamp == 0 {
		in.Timestamp = time.Now().UnixMilli()
	}
	err := l.svcCtx.AbnormalRepo.Insert(l.ctx, &deviceLog.Abnormal{
		ProductID:  in.ProductID,
		Action:     in.Action,
		Timestamp:  time.UnixMilli(in.Timestamp), // 操作时间
		DeviceName: in.DeviceName,
		TraceID:    utils.TraceIdFromContext(l.ctx),
		Reason:     in.Reason,
		Type:       in.Type,
	})
	return &dm.Empty{}, err
}
