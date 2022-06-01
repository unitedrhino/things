package logic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	mysql "github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"
	"time"

	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ManageDeviceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageDeviceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageDeviceLogic {
	return &ManageDeviceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *ManageDeviceLogic) CheckDevice(in *dm.ManageDeviceReq) (bool, error) {
	_, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(in.Info.ProductID, in.Info.DeviceName)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

/*
发现返回true 没有返回false
*/
func (l *ManageDeviceLogic) CheckProduct(in *dm.ManageDeviceReq) (bool, error) {
	_, err := l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (l *ManageDeviceLogic) AddDevice(in *dm.ManageDeviceReq) (*dm.DeviceInfo, error) {
	find, err := l.CheckDevice(in)
	if err != nil {
		l.Errorf("AddDevice|CheckDevice|in=%v\n", in)
		return nil, errors.Database.AddDetail(err.Error())
	} else if find == true {
		return nil, errors.Duplicate.AddDetail("DeviceName:" + in.Info.DeviceName)
	}
	find, err = l.CheckProduct(in)
	if err != nil {
		l.Errorf("AddDevice|CheckProduct|in=%v\n", in)
		return nil, errors.Database.AddDetail(err.Error())
	} else if find == false {
		return nil, errors.Parameter.AddDetail("not find product id:" + cast.ToString(in.Info.ProductID))
	}
	pt, err := l.svcCtx.TemplateRepo.GetTemplate(l.ctx, in.Info.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	}
	err = l.svcCtx.DeviceDataRepo.InitDevice(l.ctx, pt, in.Info.ProductID, in.Info.DeviceName)
	if err != nil {
		return nil, errors.Database.AddDetail(err.Error())
	}
	di := mysql.DeviceInfo{
		ProductID:   in.Info.ProductID,  // 产品id
		DeviceName:  in.Info.DeviceName, // 设备名称
		Secret:      utils.GetPwdBase64(20),
		Version:     in.Info.Version.GetValue(),
		CreatedTime: time.Now(),
	}
	if in.Info.LogLevel != def.UNKNOWN {
		di.LogLevel = device.LOG_CLOSE
	}
	_, err = l.svcCtx.DeviceInfo.Insert(&di)
	if err != nil {
		l.Errorf("AddDevice|DeviceInfo|Insert|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return ToDeviceInfo(&di), nil
}

func ChangeDevice(old *mysql.DeviceInfo, data *dm.DeviceInfo) {
	var isModify bool = false
	defer func() {
		if isModify {
			old.UpdatedTime = sql.NullTime{Valid: true, Time: time.Now()}
		}
	}()

	if data.LogLevel != def.UNKNOWN {
		old.LogLevel = data.LogLevel
		isModify = true
	}
	if data.Version != nil {
		old.Version = data.Version.GetValue()
		isModify = true
	}
}

func (l *ManageDeviceLogic) ModifyDevice(in *dm.ManageDeviceReq) (*dm.DeviceInfo, error) {
	di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(in.Info.ProductID, in.Info.DeviceName)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetailf("not find device|productid=%s|deviceName=%s",
				in.Info.ProductID, in.Info.DeviceName)
		}
		return nil, errors.System.AddDetail(err.Error())
	}
	ChangeDevice(di, in.Info)

	err = l.svcCtx.DeviceInfo.Update(di)
	if err != nil {
		l.Errorf("ModifyDevice|DeviceInfo|Update|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return ToDeviceInfo(di), nil
}

func (l *ManageDeviceLogic) DelDevice(in *dm.ManageDeviceReq) (*dm.DeviceInfo, error) {
	di, err := l.svcCtx.DeviceInfo.FindOneByProductIDDeviceName(in.Info.ProductID, in.Info.DeviceName)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetailf("not find device|productid=%s|deviceName=%s",
				in.Info.ProductID, in.Info.DeviceName)
		}
		l.Errorf("DelDevice|DeviceInfo|FindOne|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	{ //删除时序数据库中的表数据
		template, err := l.svcCtx.TemplateRepo.GetTemplate(l.ctx, in.Info.ProductID)
		if err != nil {
			l.Errorf("DelDevice|TemplateRepo|GetTemplate|err=%+v", err)
			return nil, errors.System.AddDetail(err.Error())
		}
		err = l.svcCtx.HubLogRepo.DropDevice(l.ctx, in.Info.ProductID, in.Info.DeviceName)
		if err != nil {
			l.Errorf("DelDevice|DeviceLogRepo|DropDevice|err=%+v", err)
			return nil, err
		}
		err = l.svcCtx.DeviceDataRepo.DropDevice(l.ctx, template, in.Info.ProductID, in.Info.DeviceName)
		if err != nil {
			l.Errorf("DelDevice|DeviceDataRepo|DropDevice|err=%+v", err)
			return nil, err
		}
	}

	err = l.svcCtx.DeviceInfo.Delete(di.Id)
	if err != nil {
		l.Errorf("DelDevice|DeviceInfo|Delete|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return &dm.DeviceInfo{}, nil
}

func (l *ManageDeviceLogic) ManageDevice(in *dm.ManageDeviceReq) (*dm.DeviceInfo, error) {
	defer func() {
		if p := recover(); p != nil {
			utils.HandleThrow(p)
		}
	}()
	l.Infof("ManageDevice|req=%+v", in)
	switch in.Opt {
	case def.OPT_ADD:
		return l.AddDevice(in)
	case def.OPT_MODIFY:
		return l.ModifyDevice(in)
	case def.OPT_DEL:
		return l.DelDevice(in)
	default:
		return nil, errors.Parameter.AddDetail("not support opt:" + string(in.Opt))
	}
}
