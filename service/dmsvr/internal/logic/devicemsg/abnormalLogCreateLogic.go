package devicemsglogic

import (
	"context"
	"gitee.com/unitedrhino/core/share/dataType"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/share/devices"
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
	di, err := l.svcCtx.DeviceCache.GetData(l.ctx, devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName})
	err = l.svcCtx.AbnormalRepo.Insert(l.ctx, &deviceLog.Abnormal{
		TenantCode: dataType.TenantCode(di.TenantCode),
		ProjectID:  dataType.ProjectID(di.ProjectID),
		AreaID:     dataType.AreaID(di.AreaID),
		AreaIDPath: dataType.AreaIDPath(di.AreaIDPath),
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
