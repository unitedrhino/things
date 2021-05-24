package logic

import (
	"context"
	"database/sql"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/spf13/cast"
	"time"

	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
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
func (l *ManageDeviceLogic) CheckDevice(in *dm.ManageDeviceReq) (bool,error){
	_,err :=l.svcCtx.DeviceInfo.FindOneByDeviceName(in.Info.DeviceName)
	switch err {
	case model.ErrNotFound:
		return false,nil
	case nil:
		return true,nil
	default:
		return false, err
	}
}
/*
发现返回true 没有返回false
*/
func (l *ManageDeviceLogic) CheckProduct(in *dm.ManageDeviceReq) (bool,error){
	_,err :=l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
	switch err {
	case model.ErrNotFound:
		return false,nil
	case nil:
		return true,nil
	default:
		return false, err
	}
}

func (l *ManageDeviceLogic) AddDevice(in *dm.ManageDeviceReq)(*dm.DeviceInfo, error){
	find,err := l.CheckDevice(in)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	}else if find == true{
		return nil, errors.Duplicate.AddDetail("DeviceName:"+in.Info.DeviceName)
	}
	l.Infof("find=%v|err=%v\n",find, err)
	find,err = l.CheckProduct(in)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	}else if find == false{
		return nil,errors.Parameter.AddDetail("not find product id:"+cast.ToString(in.Info.ProductID))
	}
	di :=  model.DeviceInfo{
		ProductID   :in.Info.ProductID,// 产品id
		DeviceID    :l.svcCtx.DeviceID.GetSnowflakeId(),// 设备id
		DeviceName  :in.Info.DeviceName,// 设备名称
		CreatedTime :time.Now(),
	}
	_,err = l.svcCtx.DeviceInfo.Insert(di)
	if err != nil {
		l.Errorf("AddDevice|DeviceInfo|Insert|err=%+v",err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return &dm.DeviceInfo{
		ProductID   :di.ProductID,     //产品id
		DeviceID    :di.DeviceID,       //设备id
		DeviceName  :di.DeviceName,    //设备名
		CreatedTime :di.CreatedTime.Unix(), //创建时间
	},nil
}

func ChangeDevice(old *model.DeviceInfo,data *dm.DeviceInfo){
	var isModify bool = false
	defer func() {
		if isModify{
			old.UpdatedTime = sql.NullTime{Valid: true,Time: time.Now()}
		}
	}()

	if data.DeviceName != "" {
		old.DeviceName = data.DeviceName
		isModify = true
	}
}

func (l *ManageDeviceLogic) ModifyDevice(in *dm.ManageDeviceReq)(*dm.DeviceInfo, error){
	di, err:= l.svcCtx.DeviceInfo.FindOne(in.Info.DeviceID)
	if err != nil {
		if err == model.ErrNotFound{
			return nil, errors.Parameter.AddDetail("not find device id:"+cast.ToString(in.Info.DeviceID))
		}
		return nil,errors.System.AddDetail(err.Error())
	}
	ChangeDevice(di,in.Info)

	err = l.svcCtx.DeviceInfo.Update(*di)
	if err != nil {
		l.Errorf("ModifyDevice|DeviceInfo|Update|err=%+v",err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return &dm.DeviceInfo{
		ProductID   :di.ProductID,     //产品id
		DeviceID    :di.DeviceID,       //设备id
		DeviceName  :di.DeviceName,    //设备名
		CreatedTime :di.CreatedTime.Unix(), //创建时间
	},nil
}

func (l *ManageDeviceLogic) DelDevice(in *dm.ManageDeviceReq)(*dm.DeviceInfo, error){
	_, err:= l.svcCtx.DeviceInfo.FindOne(in.Info.DeviceID)
	if err != nil {
		if err == model.ErrNotFound{
			return nil, errors.Parameter.AddDetail("not find device id:"+cast.ToString(in.Info.DeviceID))
		}
		l.Errorf("DelDevice|DeviceInfo|FindOne|err=%+v",err)
		return nil,errors.System.AddDetail(err.Error())
	}
	err = l.svcCtx.DeviceInfo.Delete(in.Info.DeviceID)
	if err != nil {
		l.Errorf("DelDevice|DeviceInfo|Delete|err=%+v",err)
		return nil,errors.System.AddDetail(err.Error())
	}
	return &dm.DeviceInfo{},nil
}

func (l *ManageDeviceLogic) ManageDevice(in *dm.ManageDeviceReq) (*dm.DeviceInfo, error) {
	l.Infof("ManageDevice|req=%+v",in)
	switch in.Opt {
	case dm.OptType_ADD:
		return l.AddDevice(in)
	case dm.OptType_MODIFY:
		return l.ModifyDevice(in)
	case dm.OptType_DEL:
		return l.DelDevice(in)
	default:
		return nil,errors.Parameter.AddDetail("not suppot opt:"+string(in.Opt))
	}
}
