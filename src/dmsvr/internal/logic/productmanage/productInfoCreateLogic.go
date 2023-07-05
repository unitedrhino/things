package productmanagelogic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"path"
	"regexp"

	"github.com/i-Things/things/shared/oss"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/domain/deviceAuth"
	"github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDb *relationDB.ProductInfoRepo
}

func NewProductInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoCreateLogic {
	return &ProductInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDb:   relationDB.NewProductInfoRepo(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *ProductInfoCreateLogic) CheckProduct(in *dm.ProductInfo) (bool, error) {
	_, err := l.PiDb.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductNames: []string{in.ProductName}})
	if err == nil {
		return true, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return false, nil
	}
	return false, err
}

/*
检测productid,发现返回true 没有返回false
*/
func (l *ProductInfoCreateLogic) CheckProductID(in *dm.ProductInfo) (bool, error) {
	_, err := l.PiDb.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if err == nil {
		return true, nil
	}
	if errors.Cmp(err, errors.NotFind) {
		return false, nil
	}
	return false, err
}

/*
根据用户的输入生成对应的数据库数据
*/
func (l *ProductInfoCreateLogic) ConvProductPbToPo(in *dm.ProductInfo) (*relationDB.DmProductInfo, error) {
	pi := &relationDB.DmProductInfo{
		ProductID:   in.ProductID,   // 产品id
		ProductName: in.ProductName, // 产品名称
		Desc:        in.Desc.GetValue(),
		DevStatus:   in.DevStatus.GetValue(),
		Secret:      utils.GetRandomBase64(20),
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
	if in.ProductID != "" {
		expr := `[0-9A-Za-z]{11,11}`
		match, _ := regexp.MatchString(expr, in.ProductID)
		fmt.Println(match)
		if !match {
			return nil, errors.Parameter.WithMsg("产品id格式不对,格式为11位数字和英文字母组成的字符串")
		}
		find, err := l.CheckProductID(in)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		} else if find {
			return nil, errors.Duplicate.WithMsgf("产品id重复:%s", in.ProductID).AddDetail("ProductID:" + in.ProductID)
		}
	} else {
		productID := l.svcCtx.ProductID.GetSnowflakeId() // 产品id
		in.ProductID = deviceAuth.GetStrProductID(productID)
	}
	pi, err := l.ConvProductPbToPo(in)
	if err != nil {
		return nil, err
	}

	err = l.InitProduct(pi)
	if err != nil {
		return nil, err
	}

	err = l.PiDb.Insert(l.ctx, pi)
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, err
	}
	return &dm.Response{}, nil
}
func (l *ProductInfoCreateLogic) InitProduct(pi *relationDB.DmProductInfo) error {
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
