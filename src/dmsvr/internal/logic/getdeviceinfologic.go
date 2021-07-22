package logic

import (
	"context"
	"gitee.com/godLei6/things/shared/def"
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

func (l *GetDeviceInfoLogic) GetDeviceInfo(in *dm.GetDeviceInfoReq) (resp *dm.GetDeviceInfoResp, err error) {
	l.Infof("GetDeviceInfo|req=%+v", in)
	var info []*dm.DeviceInfo
	var size int64
	if in.Page == nil || in.Page.Page == 0 {
		di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(in.ProductID, in.DeviceName)
		if err != nil {
			return nil, err
		}
		info = append(info, DBToRPCFmt(di).(*dm.DeviceInfo))
	} else {
		size, err = l.svcCtx.DmDB.GetCountByProductID(
			in.ProductID)
		if err != nil {
			return nil, err
		}
		di, err := l.svcCtx.DmDB.FindByProductID(
			in.ProductID, def.PageInfo{PageSize: in.Page.PageSize, Page: in.Page.Page})
		if err != nil {
			return nil, err
		}
		info = make([]*dm.DeviceInfo, 0, len(di))
		for _, v := range di {
			info = append(info, DBToRPCFmt(v).(*dm.DeviceInfo))
		}
	}
	return &dm.GetDeviceInfoResp{Info: info, Total: size}, nil
}
