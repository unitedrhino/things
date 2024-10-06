package devicemsglogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatusLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatusLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusLogIndexLogic {
	return &StatusLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *StatusLogIndexLogic) StatusLogIndex(in *dm.StatusLogIndexReq) (*dm.StatusLogIndexResp, error) {
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
	if err != nil {
		return nil, err
	}
	filter := deviceLog.StatusFilter{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		Status:     in.Status,
	}
	page := def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	}
	logs, err := l.svcCtx.StatusRepo.GetDeviceLog(l.ctx, filter, page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	total, err := l.svcCtx.StatusRepo.GetCountLog(l.ctx, filter, page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	var data []*dm.StatusLogInfo
	for _, v := range logs {
		data = append(data, ToDataStatusLogIndex(v))
	}
	return &dm.StatusLogIndexResp{List: data, Total: total}, nil
}
