package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataSdkLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataSdkLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataSdkLogIndexLogic {
	return &DataSdkLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备sdk调试日志
func (l *DataSdkLogIndexLogic) DataSdkLogIndex(in *dm.DataSdkLogIndexReq) (*dm.DataSdkLogIndexResp, error) {
	logs, err := l.svcCtx.SDKLogRepo.GetDeviceSDKLog(l.ctx, in.ProductID, in.DeviceName, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	var data []*dm.DataSdkLogIndex
	for _, v := range logs {
		data = append(data, ToDataSdkLogIndex(v))
	}
	//todo 总数未统计
	return &dm.DataSdkLogIndexResp{List: data}, nil
}
