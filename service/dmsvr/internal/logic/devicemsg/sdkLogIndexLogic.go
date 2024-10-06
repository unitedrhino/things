package devicemsglogic

import (
	"context"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/things/service/dmsvr/internal/domain/deviceLog"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
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
func (l *SdkLogIndexLogic) SdkLogIndex(in *dm.SdkLogIndexReq) (*dm.SdkLogIndexResp, error) {
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
	if err != nil {
		return nil, err
	}
	filter := deviceLog.SDKFilter{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
		LogLevel:   int(in.LogLevel),
	}
	logs, err := l.svcCtx.SDKLogRepo.GetDeviceSDKLog(l.ctx, filter, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	var data []*dm.SdkLogInfo
	for _, v := range logs {
		data = append(data, ToDataSdkLogIndex(v))
	}
	total, err := l.svcCtx.SDKLogRepo.GetCountLog(l.ctx, filter, def.PageInfo2{
		TimeStart: in.TimeStart,
		TimeEnd:   in.TimeEnd,
		Page:      in.Page.GetPage(),
		Size:      in.Page.GetSize(),
	})
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	//todo 总数未统计
	return &dm.SdkLogIndexResp{List: data, Total: total}, nil
}
