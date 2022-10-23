package devicemsglogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type SdkLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSdkLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SdkLogIndexLogic {
	return &SdkLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备sdk调试日志
func (l *SdkLogIndexLogic) SdkLogIndex(in *di.SdkLogIndexReq) (*di.SdkLogIndexResp, error) {
	logs, err := l.svcCtx.SDKLogRepo.GetDeviceSDKLog(l.ctx, in.ProductID, in.DeviceName, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	var data []*di.SdkLogIndex
	for _, v := range logs {
		data = append(data, ToDataSdkLogIndex(v))
	}
	//todo 总数未统计
	return &di.SdkLogIndexResp{List: data}, nil
}
