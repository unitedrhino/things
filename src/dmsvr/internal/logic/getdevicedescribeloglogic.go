package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceDescribeLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceDescribeLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceDescribeLogLogic {
	return &GetDeviceDescribeLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备调试信息记录登入登出,操作
func (l *GetDeviceDescribeLogLogic) GetDeviceDescribeLog(in *dm.GetDeviceDescribeLogReq) (*dm.GetDeviceDescribeLogResp, error) {
	logs, err := l.svcCtx.DeviceLogRepo.GetDeviceLog(l.ctx, in.ProductID, in.DeviceName, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Limit:     in.Limit,
	})
	if err != nil {
		return nil, errors.Database
	}
	var data []*dm.DeviceDescribeLog
	for _, v := range logs {
		data = append(data, ToDeviceDescribeLog(v))
	}

	return &dm.GetDeviceDescribeLogResp{List: data}, nil
}
