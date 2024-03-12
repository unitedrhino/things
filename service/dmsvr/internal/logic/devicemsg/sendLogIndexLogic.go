package devicemsglogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogIndexLogic {
	return &SendLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendLogIndexLogic) SendLogIndex(in *dm.SendLogIndexReq) (*dm.SendLogIndexResp, error) {
	filter := deviceLog.SendFilter{
		ProjectID:  in.ProjectID,
		AreaIDs:    in.AreaIDs,
		UserID:     in.UserID,
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		Actions:    in.Actions,
		ResultCode: in.ResultCode,
	}
	page := def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	}
	logs, err := l.svcCtx.SendRepo.GetDeviceLog(l.ctx, filter, page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	total, err := l.svcCtx.SendRepo.GetCountLog(l.ctx, filter, page)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	var data []*dm.SendLogInfo
	for _, v := range logs {
		data = append(data, ToDataSendLogIndex(v))
	}
	return &dm.SendLogIndexResp{List: data, Total: total}, nil

}
