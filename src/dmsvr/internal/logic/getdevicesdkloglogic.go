package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceSDKLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceSDKLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceSDKLogLogic {
	return &GetDeviceSDKLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备sdk调试日志
func (l *GetDeviceSDKLogLogic) GetDeviceSDKLog(in *dm.GetDeviceSDKLogReq) (*dm.GetDeviceSDKLogResp, error) {
	var page int64 = 1
	var pageSize int64 = 20
	if !(in.Page == nil || in.Page.Page == 0 || in.Page.PageSize == 0) {
		page = in.Page.Page
		pageSize = in.Page.PageSize
	}
	logs, err := l.svcCtx.SDKLogRepo.GetDeviceSDKLog(l.ctx, in.ProductID, in.DeviceName, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Limit:     in.Limit,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, errors.Database
	}
	var data []*dm.DeviceSDKLog
	for _, v := range logs {
		data = append(data, ToDeviceSDKLog(v))
	}
	return &dm.GetDeviceSDKLogResp{List: data}, nil
}
