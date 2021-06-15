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

func (l *GetDeviceInfoLogic) GetDeviceInfo(in *dm.GetDeviceInfoReq) (*dm.GetDeviceInfoResp, error) {
	l.Infof("GetDeviceInfo|req=%+v",in)
	di,err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(in.ProductID,in.DeviceName)
	if err != nil {
		return nil, err
	}
	return &dm.GetDeviceInfoResp{Info: []*dm.DeviceInfo{{
		ProductID:di.ProductID,
		DeviceName: di.DeviceName,
		CreatedTime: di.CreatedTime.Unix(),
		FirstLogin: di.FirstLogin.Time.Unix(),
		LastLogin: di.LastLogin.Time.Unix(),
		Secret: di.Secret,
	}}}, nil
}
