package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataHubLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataHubLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataHubLogIndexLogic {
	return &DataHubLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备调试信息记录登入登出,操作
func (l *DataHubLogIndexLogic) DataHubLogIndex(in *dm.DataHubLogIndexReq) (*dm.DataHubLogIndexResp, error) {
	logs, err := l.svcCtx.HubLogRepo.GetDeviceLog(l.ctx, in.ProductID, in.DeviceName, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	total, err := l.svcCtx.HubLogRepo.GetCountLog(l.ctx, in.ProductID, in.DeviceName, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	var data []*dm.DataHubLogIndex
	for _, v := range logs {
		data = append(data, ToDataHubLogIndex(v))
	}
	return &dm.DataHubLogIndexResp{List: data, Total: total}, nil
}
