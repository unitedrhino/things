package devicemsglogic

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
	if err != nil {
		return nil, err
	}
	filter := deviceLog.SendFilter{
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
