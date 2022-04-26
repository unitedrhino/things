package logic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"
	"github.com/i-Things/things/src/dmsvr/internal/domain/device"
	"github.com/i-Things/things/src/dmsvr/internal/domain/productDetail"
	"github.com/i-Things/things/src/dmsvr/internal/domain/templateModel"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/spf13/cast"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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
	pi, pt := l.InsertProduct(in)
	t, _ := templateModel.NewTemplate([]byte(pt.Template))
	if err := l.svcCtx.DeviceDataRepo.InitProduct(
		l.ctx, t, pi.ProductID); err != nil {
		l.Errorf("%s InitProduct failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.DmDB.Insert(pi, pt)
	if err != nil {
		l.Errorf("AddProduct|Insert|err=%+v", err)
		return nil, errors.System.AddDetail(err.Error())
	}
	return ToProductInfo(pi), nil
}

/*
根据用户的输入生成对应的数据库数据
*/
func (l *ManageProductLogic) InsertProduct(in *dm.ManageProductReq) (*mysql.ProductInfo, *mysql.ProductTemplate) {
	info := in.Info
	ProductID := l.svcCtx.ProductID.GetSnowflakeId() // 产品id
	createTime := time.Now()
	pt := &mysql.ProductTemplate{
		Template:    templateModel.DefaultTemplate,
		ProductID:   device.GetStrProductID(ProductID),
		CreatedTime: createTime,
	}
	pi := &mysql.ProductInfo{
		ProductID:   device.GetStrProductID(ProductID), // 产品id
		ProductName: info.ProductName,                  // 产品名称
		Description: info.Description.GetValue(),
		DevStatus:   info.DevStatus.GetValue(),
		CreatedTime: createTime,
	}
	if info.AutoRegister != def.UNKNOWN {
		pi.AutoRegister = info.AutoRegister
	} else {
		pi.AutoRegister = productDetail.AUTO_REG_CLOSE
	}
	if info.DataProto != def.UNKNOWN {
		pi.DataProto = info.DataProto
	} else {
		pi.DataProto = productDetail.DATA_CUSTOM
	}
	if info.DeviceType != def.UNKNOWN {
		pi.DeviceType = info.DeviceType
	} else {
		pi.DeviceType = productDetail.DEV_DEVICE
	}
	if info.NetType != def.UNKNOWN {
		pi.NetType = info.NetType
	} else {
		pi.NetType = productDetail.NET_OTHER
	}
	if info.DeviceType != def.UNKNOWN {
		pi.DeviceType = info.DeviceType
	} else {
		pi.DeviceType = productDetail.DEV_DEVICE
	}
	if info.AuthMode != def.UNKNOWN {
		pi.AuthMode = info.AuthMode
	} else {
		pi.AuthMode = productDetail.AUTH_PWD
	}
	return pi, pt
}

func UpdateProductInfo(old *mysql.ProductInfo, data *dm.ProductInfo) {
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

	if data.AutoRegister != def.UNKNOWN {
		old.AutoRegister = data.AutoRegister
		isModify = true
	}
	if data.DevStatus != nil {
		old.DevStatus = data.DevStatus.GetValue()
		isModify = true
	}

	if data.ProductName != "" {
		old.ProductName = data.ProductName
		isModify = true
	}
	if data.AuthMode != 0 {
		old.AuthMode = data.AuthMode
		isModify = true
	}
	if data.DeviceType != 0 {
		old.DeviceType = data.DeviceType
		isModify = true
	}
	if data.CategoryID != 0 {
		old.CategoryID = data.CategoryID
		isModify = true
	}
	if data.NetType != 0 {
		old.NetType = data.NetType
		isModify = true
	}
	if data.DataProto != 0 {
		old.DataProto = data.DataProto
		isModify = true
	}
	if data.AutoRegister != 0 {
		old.AutoRegister = data.AutoRegister
		isModify = true
	}

}

func (l *ManageProductLogic) ModifyProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	pi, err := l.svcCtx.ProductInfo.FindOne(in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.Database.AddDetail(err.Error())
	}
	UpdateProductInfo(pi, in.Info)

	err = l.svcCtx.ProductInfo.Update(pi)
	if err != nil {
		l.Errorf("ModifyProduct|ProductInfo|Update|err=%+v", err)
		return nil, errors.Database.AddDetail(err.Error())
	}
	return ToProductInfo(pi), nil
}

func (l *ManageProductLogic) DelProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	pt, err := l.svcCtx.TemplateRepo.GetTemplate(l.ctx, in.Info.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err.Error())
	}
	err = l.svcCtx.DeviceDataRepo.DropProduct(l.ctx, pt, in.Info.ProductID)
	if err != nil {
		l.Errorf("DelProduct|DropProduct|err=%+v", err)
		return nil, errors.Database.AddDetail(err.Error())
	}
	l.svcCtx.TemplateRepo.ClearCache(l.ctx, in.Info.ProductID)
	err = l.svcCtx.DmDB.Delete(in.Info.ProductID)
	if err != nil {
		l.Errorf("DelProduct|Delete|err=%+v", err)
		return nil, errors.Database.AddDetail(err.Error())
	}
	err = l.svcCtx.DataUpdate.TempModelUpdate(l.ctx, &templateModel.TemplateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.ProductInfo{}, nil
}

func (l *ManageProductLogic) ManageProduct(in *dm.ManageProductReq) (*dm.ProductInfo, error) {
	l.Infof("ManageProduct|opt=%d|info=%+v", in.Opt, in.Info)
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
