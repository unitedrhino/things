package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeviceInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoDeleteLogic {
	return &DeviceInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除设备
func (l *DeviceInfoDeleteLogic) DeviceInfoDelete(in *dm.DeviceInfoDeleteReq) (*dm.Response, error) {
	di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(l.ctx, in.ProductID, in.DeviceName)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetailf("not find device|productid=%s|deviceName=%s",
				in.ProductID, in.DeviceName)
		}
		l.Errorf("DelDevice|DeviceInfo|FindOne|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	{ //删除时序数据库中的表数据
		template, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
		if err != nil {
			l.Errorf("DelDevice|SchemaRepo|GetSchemaModel|err=%+v", err)
			return nil, errors.System.AddDetail(err)
		}
		err = l.svcCtx.HubLogRepo.DropDevice(l.ctx, in.ProductID, in.DeviceName)
		if err != nil {
			l.Errorf("DelDevice|DeviceLogRepo|DropDevice|err=%+v", err)
			return nil, err
		}
		err = l.svcCtx.DeviceDataRepo.DropDevice(l.ctx, template, in.ProductID, in.DeviceName)
		if err != nil {
			l.Errorf("DelDevice|DeviceDataRepo|DropDevice|err=%+v", err)
			return nil, err
		}
	}

	err = l.svcCtx.DeviceInfo.Delete(l.ctx, di.Id)
	if err != nil {
		l.Errorf("DelDevice|DeviceInfo|Delete|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.Response{}, nil
}
