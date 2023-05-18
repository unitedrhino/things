package productmanagelogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/oss"
	"path"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
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
func (l *ProductInfoCreateLogic) ConvProductPbToPo(in *dm.ProductInfo) (*mysql.DmProductInfo, error) {
	ProductID := l.svcCtx.ProductID.GetSnowflakeId() // 产品id
	pi := &mysql.DmProductInfo{
		ProductID:   deviceAuth.GetStrProductID(ProductID), // 产品id
		ProductName: in.ProductName,                        // 产品名称
		Desc:        in.Desc.GetValue(),
		DevStatus:   in.DevStatus.GetValue(),
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
	if in.Tags != nil {
		tags, err := json.Marshal(in.Tags)
		if err == nil {
			pi.Tags = string(tags)
		}
	} else {
		pi.Tags = "{}"
	}
	if in.ProductImg != "" { //如果填了参数且不等于原来的,说明修改头像,需要处理
		si, err := oss.GetSceneInfo(in.ProductImg)
		if err != nil {
			return nil, err
		}
		if !(si.Business == oss.BusinessProductManage && si.Scene == oss.SceneProductImg) {
			return nil, errors.Parameter.WithMsg("产品图片的路径不对")
		}
		si.FilePath = pi.ProductID + path.Ext(si.FilePath)
		nwePath, err := oss.GetFilePath(si, false)
		if err != nil {
			return nil, err
		}
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(in.ProductImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		pi.ProductImg = path
	}
	return pi, nil
}

// 新增设备
func (l *ProductInfoCreateLogic) ProductInfoCreate(in *dm.ProductInfo) (*dm.Response, error) {
	find, err := l.CheckProduct(in)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.WithMsgf("产品名称重复:%s", in.ProductName).AddDetail("ProductName:" + in.ProductName)
	}

	pi, err := l.ConvProductPbToPo(in)
	if err != nil {
		return nil, err
	}

	err = l.InitProduct(pi)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.ProductInfo.Insert(l.ctx, pi)
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}

	return &dm.Response{}, nil
}
func (l *ProductInfoCreateLogic) InitProduct(pi *mysql.DmProductInfo) error {
	t, _ := schema.NewSchemaTsl([]byte(schema.DefaultSchema))
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
