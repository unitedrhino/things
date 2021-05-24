package logic

import (
	"context"

	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetDeviceInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceInfoLogic {
	return &GetDeviceInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceInfoLogic) GetDeviceInfo(in *dm.GetDeviceInfoReq) (*dm.DeviceInfo, error) {
	l.Infof("GetDeviceInfo|req=%+v",in)
	di,err := l.svcCtx.DeviceInfo.FindOne(in.DeviceID)
	if err != nil {
		return nil, err
	}
	return &dm.DeviceInfo{
		ProductID:di.ProductID,
		DeviceID: di.DeviceID,
		DeviceName: di.DeviceName,
		CreatedTime: di.CreatedTime.Unix(),
	}, nil
}
