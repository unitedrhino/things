package logic

import (
	"context"
	"database/sql"
	"github.com/go-things/things/shared/def"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/src/dmsvr/internal/repo/model/mysql"
	"github.com/spf13/cast"
	"time"

	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"

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
func (l *ManageProductLogic) CheckProduct(in *dm.ManageProductReq) (bool, error) {
	_, err := l.svcCtx.ProductInfo.FindOneByProductName(in.Info.ProductName)
	switch err {
	case mysql.ErrNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}

func (l *ManageProductLogic) AddProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	find, err := l.CheckProduct(in)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	} else if find == true {
		return nil, errors.Duplicate.AddDetail("ProductName:" + in.Info.ProductName)
	}
	pi := l.InsertProduct(in)
	_, err = l.svcCtx.ProductInfo.Insert(*pi)
	if err != nil {
		l.Errorf("AddProduct|ProductInfo|Insert|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return DBToRPCFmt(pi).(*dm.ProductInfo), nil
}

/*
根据用户的输入生成对应的数据库数据
*/
func (l *ManageProductLogic) InsertProduct(in *dm.ManageProductReq) *mysql.ProductInfo {
	info := in.Info
	ProductID := l.svcCtx.ProductID.GetSnowflakeId() // 产品id
	pi := &mysql.ProductInfo{
		ProductID:   dm.GetStrProductID(ProductID), // 产品id
		ProductName: info.ProductName,              // 产品名称
		Description: info.Description.GetValue(),
		Template:    info.Template.GetValue(),
		DevStatus:   info.DevStatus.GetValue(),
		CreatedTime: time.Now(),
	}
	if info.AutoRegister != def.UNKNOWN {
		pi.AutoRegister = info.AutoRegister
	} else {
		pi.AutoRegister = dm.AUTO_REG_CLOSE
	}
	if info.DataProto != def.UNKNOWN {
		pi.DataProto = info.DataProto
	} else {
		pi.DataProto = dm.DATA_CUSTOM
	}
	if info.DeviceType != def.UNKNOWN {
		pi.DeviceType = info.DeviceType
	} else {
		pi.DeviceType = dm.DEV_DEVICE
	}
	if info.NetType != def.UNKNOWN {
		pi.NetType = info.NetType
	} else {
		pi.NetType = dm.NET_OTHER
	}
	if info.DeviceType != def.UNKNOWN {
		pi.DeviceType = info.DeviceType
	} else {
		pi.DeviceType = dm.DEV_DEVICE
	}
	if info.AuthMode != def.UNKNOWN {
		pi.AuthMode = info.AuthMode
	} else {
		pi.AuthMode = dm.AUTH_PWD
	}
	return pi
}

func UpdateProduct(old *mysql.ProductInfo, data *dm.ProductInfo) {
	var isModify bool = false
	defer func() {
		if isModify {
			old.UpdatedTime = sql.NullTime{Valid: true, Time: time.Now()}
		}
	}()
	if data.ProductName != "" {
		old.ProductName = data.ProductName
		isModify = true
	}
	if data.AuthMode != def.UNKNOWN {
		old.AuthMode = data.AuthMode
		isModify = true
	}
	if data.Description != nil {
		old.Description = data.Description.GetValue()
		isModify = true
	}
	if data.Template != nil {
		old.Template = data.Template.GetValue()
		isModify = true
	}
	if data.AutoRegister != def.UNKNOWN {
		old.AutoRegister = data.AutoRegister
		isModify = true
	}
	if data.DevStatus != nil {
		old.DevStatus = data.DevStatus.GetValue()
		isModify = true
	}
}

func (l *ManageProductLogic) ModifyProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	pi, err := l.svcCtx.ProductInfo.FindOneByProductID(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.System.AddDetail(err.Error())
	}
	UpdateProduct(pi, in.Info)

	err = l.svcCtx.ProductInfo.Update(*pi)
	if err != nil {
		l.Errorf("ModifyProduct|ProductInfo|Update|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return DBToRPCFmt(pi).(*dm.ProductInfo), nil
}

func (l *ManageProductLogic) DelProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	info, err := l.svcCtx.ProductInfo.FindOneByProductID(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find device id:" + cast.ToString(in.Info.ProductID))
		}
		l.Errorf("DelProduct|ProductInfo|FindOne|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	err = l.svcCtx.ProductInfo.Delete(info.Id)
	if err != nil {
		l.Errorf("DelProduct|ProductInfo|Delete|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return &dm.ProductInfo{}, nil
}

func (l *ManageProductLogic) ManageProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	l.Infof("ManageProduct|opt=%d|req=%+v", in.Opt, in)
	switch in.Opt {
	case def.OPT_ADD:
		if in.Info == nil {
			return nil, errors.Parameter.WithMsg("add opt need info")
		}
		return l.AddProduct(in)
	case def.OPT_MODIFY:
		return l.ModifyProduct(in)
	case def.OPT_DEL:
		return l.DelProduct(in)
	default:
		return nil, errors.Parameter.AddDetail("not suppot opt:" + cast.ToString(in.Opt))
	}
}
