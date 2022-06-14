package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
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
	if (in.Page == nil || in.Page.Page == 0) && in.DeviceName != "" {
		di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
		if err != nil {
			if err == mysql.ErrNotFound {
				return nil, errors.NotFind
			}
			return nil, err
		}
		info = append(info, ToDeviceInfo(di))
	} else {
		var page int64 = 1
		var pageSize int64 = 20
		if !(in.Page == nil || in.Page.Page == 0 || in.Page.PageSize == 0) {
			page = in.Page.Page
			pageSize = in.Page.PageSize
		}
		size, err = l.svcCtx.DmDB.GetCountByProductID(
			l.ctx, in.ProductID)
		if err != nil {
			return nil, err
		}
		di, err := l.svcCtx.DmDB.FindByProductID(
			l.ctx, in.ProductID, def.PageInfo{PageSize: pageSize, Page: page})
		if err != nil {
			return nil, err
		}
		info = make([]*dm.DeviceInfo, 0, len(di))
		for _, v := range di {
			info = append(info, ToDeviceInfo(v))
		}
	}
	return &dm.GetDeviceInfoResp{List: info, Total: size}, nil
}
