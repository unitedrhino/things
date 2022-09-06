package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
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
		UpdatedTime: createTime,
	}
	pi := &mysql.ProductInfo{
		ProductID:   deviceAuth.GetStrProductID(ProductID), // 产品id
		ProductName: in.ProductName,                        // 产品名称
		Description: in.Description.GetValue(),
		DevStatus:   in.DevStatus.GetValue(),
		CreatedTime: createTime,
		UpdatedTime: createTime,
	}
	if in.AutoRegister != def.Unknown {
		pi.AutoRegister = in.AutoRegister
	} else {
		pi.AutoRegister = def.AutoRegClose
	}
	if in.DataProto != def.Unknown {
		pi.DataProto = in.DataProto
	} else {
		pi.DataProto = def.DataProtoCustom
	}
	if in.DeviceType != def.Unknown {
		pi.DeviceType = in.DeviceType
	} else {
		pi.DeviceType = def.DeviceTypeDevice
	}
	if in.NetType != def.Unknown {
		pi.NetType = in.NetType
	} else {
		pi.NetType = def.NetOther
	}
	if in.DeviceType != def.Unknown {
		pi.DeviceType = in.DeviceType
	} else {
		pi.DeviceType = def.DeviceTypeDevice
	}
	if in.AuthMode != def.Unknown {
		pi.AuthMode = in.AuthMode
	} else {
		pi.AuthMode = def.AuthModePwd
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
	err = l.InitProduct(pi, pt)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DmDB.Insert(l.ctx, pi, pt)
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	return &dm.Response{}, nil
}
func (l *ProductInfoCreateLogic) InitProduct(pi *mysql.ProductInfo, pt *mysql.ProductSchema) error {
	t, _ := schema.NewSchema([]byte(pt.Schema))
	if err := l.svcCtx.SchemaManaRepo.InitProduct(
		l.ctx, t, pi.ProductID); err != nil {
		l.Errorf("%s.SchemaManaRepo.InitProduct failure,err:%v", utils.FuncName(), err)
		return errors.Database.AddDetail(err)
	}
	if err := l.svcCtx.HubLogRepo.InitProduct(
		l.ctx, pi.ProductID); err != nil {
		l.Errorf("%s.HubLogRepo.InitProduct failure,err:%v", utils.FuncName(), err)
		return errors.Database.AddDetail(err)
	}
	if err := l.svcCtx.SDKLogRepo.InitProduct(
		l.ctx, pi.ProductID); err != nil {
		l.Errorf("%s.SDKLogRepo.InitProduct failure,err:%v", utils.FuncName(), err)
		return errors.Database.AddDetail(err)
	}
	return nil
}
