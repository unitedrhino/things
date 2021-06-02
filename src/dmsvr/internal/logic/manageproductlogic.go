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
	pi :=  model.ProductInfo{
		ProductID    :l.svcCtx.ProductID.GetSnowflakeId(),// 产品id
		ProductName  :in.Info.ProductName,// 产品名称
		CreatedTime :time.Now(),
	}
	_,err = l.svcCtx.ProductInfo.Insert(pi)
	if err != nil {
		l.Errorf("AddProduct|ProductInfo|Insert|err=%+v",err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return &dm.ProductInfo{
		ProductID    :pi.ProductID,     //产品id
		ProductName  :pi.ProductName,    //产品名
		CreatedTime  :pi.CreatedTime.Unix(), //创建时间
	},nil
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
	if data.AutoRegister != dm.UNKNOWN {
		old.AutoRegister = int64(data.AutoRegister)
		isModify = true
	}
}

func (l *ManageProductLogic) ModifyProduct(in *dm.ManageProductReq)(*dm.ProductInfo, error){
	pi, err:= l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
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
	return &dm.ProductInfo{
		ProductID   :pi.ProductID,     //产品id
		ProductName :pi.ProductName,       //产品名
		CreatedTime :pi.CreatedTime.Unix(), //创建时间
	},nil
}

func (l *ManageProductLogic) DelProduct(in *dm.ManageProductReq)(*dm.ProductInfo, error){
	_, err:= l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
	if err != nil {
		if err == model.ErrNotFound{
			return nil, errors.Parameter.AddDetail("not find device id:"+cast.ToString(in.Info.ProductID))
		}
		l.Errorf("DelProduct|ProductInfo|FindOne|err=%+v",err)
		return nil,errors.System.AddDetail(err.Error())
	}
	err = l.svcCtx.ProductInfo.Delete(in.Info.ProductID)
	if err != nil {
		l.Errorf("DelProduct|ProductInfo|Delete|err=%+v",err)
		return nil,errors.System.AddDetail(err.Error())
	}
	return &dm.ProductInfo{},nil
}

func (l *ManageProductLogic) ManageProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	l.Infof("ManageProduct|req=%+v",in)
	switch in.Opt {
	case dm.OPT_ADD:
		return l.AddProduct(in)
	case dm.OPT_MODIFY:
		return l.ModifyProduct(in)
	case dm.OPT_DEL:
		return l.DelProduct(in)
	default:
		return nil,errors.Parameter.AddDetail("not suppot opt:"+string(in.Opt))
	}
}
