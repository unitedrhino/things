package logic

import (
	"context"
	"database/sql"
	"gitee.com/godLei6/things/shared/errors"
	"gitee.com/godLei6/things/src/dmsvr/model"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/spf13/cast"
	"time"

	"gitee.com/godLei6/things/src/dmsvr/dm"
	"gitee.com/godLei6/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewManageProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ManageProductLogic {
	return &ManageProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}


/*
发现返回true 没有返回false
*/
func (l *ManageProductLogic) CheckProduct(in *dm.ManageProductReq) (bool,error){
	_,err :=l.svcCtx.ProductInfo.FindOneByProductName(in.Info.ProductName)
	switch err {
	case model.ErrNotFound:
		return false, nil
	case nil:
		return true,nil
	default:
		return false, err
	}
}

func (l *ManageProductLogic) AddProduct(in *dm.ManageProductReq)(*dm.ProductInfo, error){

	find,err := l.CheckProduct(in)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	}else if find == true{
		return nil,errors.Duplicate.AddDetail("ProductName:" + in.Info.ProductName)
	}
	pi :=  l.InsertProduct(in)
	_,err = l.svcCtx.ProductInfo.Insert(*pi)
	if err != nil {
		l.Errorf("AddProduct|ProductInfo|Insert|err=%+v",err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return DbToProto(pi),nil
}

func DbToProto(pi *model.ProductInfo)*dm.ProductInfo{
	return &dm.ProductInfo{
		ProductID    :pi.ProductID,     //产品id
		ProductName  :pi.ProductName,    //产品名
		AuthMode     :pi.AuthMode,//认证方式:0:账密认证,1:秘钥认证
		DeviceType   :pi.DeviceType,//设备类型:0:设备,1:网关,2:子设备
		CategoryID   :pi.CategoryID,//产品品类
		NetType      :pi.NetType,//通讯方式:0:其他,1:wi-fi,2:2G/3G/4G,3:5G,4:BLE,5:LoRaWAN
		DataProto    :pi.DataProto,//数据协议:0:自定义,1:数据模板
		AutoRegister :pi.AutoRegister,//动态注册:0:关闭,1:打开,2:打开并自动创建设备
		Secret       :pi.Secret,//动态注册产品秘钥 只读
		Description  :&wrappers.StringValue{Value: pi.Description},   //描述
		CreatedTime  :pi.CreatedTime.Unix(), //创建时间
	}
}

/*
根据用户的输入生成对应的数据库数据
*/
func (l *ManageProductLogic)InsertProduct(in *dm.ManageProductReq)(*model.ProductInfo){
	info := in.Info
	ProductID := l.svcCtx.ProductID.GetSnowflakeId()// 产品id
	pi :=  &model.ProductInfo{
		ProductID    :dm.GetStrProductID(ProductID),// 产品id
		ProductName  :info.ProductName,// 产品名称
		Description	 : info.Description.GetValue(),
		Template: info.Template.GetValue(),
		CreatedTime :time.Now(),
	}
	if info.AutoRegister != dm.UNKNOWN {
		pi.AutoRegister = info.AutoRegister
	}else {
		pi.AutoRegister = dm.AUTO_REG_CLOSE
	}
	if info.DataProto != dm.UNKNOWN {
		pi.DataProto = info.DataProto
	}else {
		pi.DataProto = dm.DATA_CUSTOM
	}
	if info.DeviceType != dm.UNKNOWN {
		pi.DeviceType = info.DeviceType
	}else {
		pi.DeviceType = dm.DEV_DEVICE
	}
	if info.NetType != dm.UNKNOWN {
		pi.NetType = info.NetType
	}else {
		pi.NetType = dm.NET_OTHER
	}
	if info.DeviceType != dm.UNKNOWN {
		pi.DeviceType = info.DeviceType
	}else {
		pi.DeviceType = dm.DEV_DEVICE
	}
	if info.AuthMode != dm.UNKNOWN {
		pi.AuthMode = info.AuthMode
	}else {
		pi.AuthMode = dm.AUTH_PWD
	}
	return pi
}

func UpdateProduct(old *model.ProductInfo,data *dm.ProductInfo){
	var isModify bool = false
	defer func() {
		if isModify{
			old.UpdatedTime = sql.NullTime{Valid: true,Time: time.Now()}
		}
	}()
	if data.ProductName != "" {
		old.ProductName = data.ProductName
		isModify = true
	}
	if data.AuthMode != dm.UNKNOWN {
		old.AuthMode = int64(data.AuthMode)
		isModify = true
	}
	if data.Description != nil{
		old.Description = data.Description.GetValue()
		isModify = true
	}
	if data.Template != nil{
		old.Template = data.Template.GetValue()
		isModify = true
	}
	if data.AutoRegister != dm.UNKNOWN {
		old.AutoRegister = int64(data.AutoRegister)
		isModify = true
	}
}

func (l *ManageProductLogic) ModifyProduct(in *dm.ManageProductReq)(*dm.ProductInfo, error){
	pi, err:= l.svcCtx.ProductInfo.FindOneByProductID(in.Info.ProductID)
	if err != nil {
		if err == model.ErrNotFound{
			return nil, errors.Parameter.AddDetail("not find ProductID id:"+cast.ToString(in.Info.ProductID))
		}
		return nil,errors.System.AddDetail(err.Error())
	}
	UpdateProduct(pi,in.Info)

	err = l.svcCtx.ProductInfo.Update(*pi)
	if err != nil {
		l.Errorf("ModifyProduct|ProductInfo|Update|err=%+v",err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return DbToProto(pi),nil
}

func (l *ManageProductLogic) DelProduct(in *dm.ManageProductReq)(*dm.ProductInfo, error){
	info, err:= l.svcCtx.ProductInfo.FindOneByProductID(in.Info.ProductID)
	if err != nil {
		if err == model.ErrNotFound{
			return nil, errors.Parameter.AddDetail("not find device id:"+cast.ToString(in.Info.ProductID))
		}
		l.Errorf("DelProduct|ProductInfo|FindOne|err=%+v",err)
		return nil,errors.System.AddDetail(err.Error())
	}
	err = l.svcCtx.ProductInfo.Delete(info.Id)
	if err != nil {
		l.Errorf("DelProduct|ProductInfo|Delete|err=%+v",err)
		return nil,errors.System.AddDetail(err.Error())
	}
	return &dm.ProductInfo{},nil
}

func (l *ManageProductLogic) ManageProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	l.Infof("ManageProduct|opt=%d|req=%+v",in.Opt,in)
	switch in.Opt {
	case dm.OPT_ADD:
		if in.Info == nil {
			return nil,errors.Parameter.WithMsg("add opt need info")
		}
		return l.AddProduct(in)
	case dm.OPT_MODIFY:
		return l.ModifyProduct(in)
	case dm.OPT_DEL:
		return l.DelProduct(in)
	default:
		return nil,errors.Parameter.AddDetail("not suppot opt:"+string(in.Opt))
	}
}
