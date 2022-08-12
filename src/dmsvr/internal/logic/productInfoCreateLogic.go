package logic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/domain/deviceAuth"
	"github.com/i-Things/things/src/dmsvr/internal/domain/productInfo"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"time"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoCreateLogic {
	return &ProductInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *ProductInfoCreateLogic) CheckProduct(in *dm.ProductInfo) (bool, error) {
	_, err := l.svcCtx.ProductInfo.FindOneByProductName(l.ctx, in.ProductName)
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
根据用户的输入生成对应的数据库数据
*/
func (l *ProductInfoCreateLogic) InsertProduct(in *dm.ProductInfo) (*mysql.ProductInfo, *mysql.ProductSchema) {
	ProductID := l.svcCtx.ProductID.GetSnowflakeId() // 产品id
	createTime := time.Now()
	pt := &mysql.ProductSchema{
		Schema:      schema.DefaultSchema,
		ProductID:   deviceAuth.GetStrProductID(ProductID),
		CreatedTime: createTime,
	}
	pi := &mysql.ProductInfo{
		ProductID:   deviceAuth.GetStrProductID(ProductID), // 产品id
		ProductName: in.ProductName,                        // 产品名称
		Description: in.Description.GetValue(),
		DevStatus:   in.DevStatus.GetValue(),
		CreatedTime: createTime,
	}
	if in.AutoRegister != def.UNKNOWN {
		pi.AutoRegister = in.AutoRegister
	} else {
		pi.AutoRegister = productInfo.AutoRegClose
	}
	if in.DataProto != def.UNKNOWN {
		pi.DataProto = in.DataProto
	} else {
		pi.DataProto = productInfo.DataCustom
	}
	if in.DeviceType != def.UNKNOWN {
		pi.DeviceType = in.DeviceType
	} else {
		pi.DeviceType = productInfo.DevDevice
	}
	if in.NetType != def.UNKNOWN {
		pi.NetType = in.NetType
	} else {
		pi.NetType = productInfo.NetOther
	}
	if in.DeviceType != def.UNKNOWN {
		pi.DeviceType = in.DeviceType
	} else {
		pi.DeviceType = productInfo.DevDevice
	}
	if in.AuthMode != def.UNKNOWN {
		pi.AuthMode = in.AuthMode
	} else {
		pi.AuthMode = productInfo.AuthPwd
	}
	return pi, pt
}

// 新增设备
func (l *ProductInfoCreateLogic) ProductInfoCreate(in *dm.ProductInfo) (*dm.Response, error) {
	find, err := l.CheckProduct(in)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.AddDetail("ProductName:" + in.ProductName)
	}
	pi, pt := l.InsertProduct(in)
	t, _ := schema.NewSchema([]byte(pt.Schema))
	if err := l.svcCtx.HubLogRepo.InitProduct(
		l.ctx, pi.ProductID); err != nil {
		l.Errorf("%s|DeviceLogRepo|InitProduct| failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	if err := l.svcCtx.DeviceDataRepo.InitProduct(
		l.ctx, t, pi.ProductID); err != nil {
		l.Errorf("%s|DeviceDataRepo|InitProduct| failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.DmDB.Insert(l.ctx, pi, pt)
	if err != nil {
		l.Errorf("AddProduct|Insert|err=%+v", err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.Response{}, nil
}
