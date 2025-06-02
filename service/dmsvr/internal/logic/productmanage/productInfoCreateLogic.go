package productmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/protocol"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/unitedrhino/things/share/topics"
	"regexp"

	"gitee.com/unitedrhino/share/oss"

	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/domain/deviceAuth"
	"gitee.com/unitedrhino/things/share/domain/schema"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
}

func NewProductInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoCreateLogic {
	return &ProductInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
	}
}

/*
发现返回true 没有返回false
*/
func (l *ProductInfoCreateLogic) CheckProduct(in *dm.ProductInfo) (bool, error) {
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductNames: []string{in.ProductName}})
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
	_, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
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
		ProductID:        in.ProductID,   // 产品id
		ProductName:      in.ProductName, // 产品名称
		Desc:             in.Desc.GetValue(),
		Status:           in.Status,
		Secret:           utils.GetRandomBase64(20),
		TrialTime:        in.TrialTime.GetValue(),
		SceneMode:        in.SceneMode,
		DeviceSchemaMode: in.DeviceSchemaMode,
		BindLevel:        in.BindLevel,
		SubProtocolCode:  in.SubProtocolCode.GetValue(),
	}
	if in.AutoRegister != def.Unknown {
		pi.AutoRegister = in.AutoRegister
	}
	if in.ProtocolCode != "" {
		pi.ProtocolCode = in.ProtocolCode
	} else {
		pi.ProtocolCode = protocol.CodeIThings
	}
	if in.DeviceType != def.Unknown {
		pi.DeviceType = in.DeviceType
	}
	if in.CategoryID != 0 {
		pi.CategoryID = in.CategoryID
	}
	if in.NetType != def.Unknown {
		pi.NetType = in.NetType
	}
	if in.AuthMode != def.Unknown {
		pi.AuthMode = in.AuthMode
	}
	if in.Tags == nil {
		in.Tags = map[string]string{}
	}
	if in.ProtocolConf == nil {
		in.ProtocolConf = map[string]string{}
	}
	if in.SubProtocolConf == nil {
		in.SubProtocolConf = map[string]string{}
	}
	pi.Tags = in.Tags
	if in.ProductImg != "" { //如果填了参数且不等于原来的,说明修改头像,需要处理
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProductManage, oss.SceneProductImg, fmt.Sprintf("%s/%s", in.ProductID, oss.GetFileNameWithPath(in.ProductImg)))
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(in.ProductImg, nwePath)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		}
		pi.ProductImg = path
	}
	return pi, nil
}

// 新增设备
func (l *ProductInfoCreateLogic) ProductInfoCreate(in *dm.ProductInfo) (*dm.Empty, error) {
	if err := ctxs.IsRoot(l.ctx); err != nil {
		return nil, err
	}
	find, err := l.CheckProduct(in)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	} else if find == true {
		return nil, errors.Duplicate.WithMsgf("产品名称重复:%s", in.ProductName).AddDetail("ProductName:" + in.ProductName)
	}
	if in.ProductID != "" {
		expr := `[0-9A-Za-z]{2,20}`
		match, _ := regexp.MatchString(expr, in.ProductID)
		//fmt.Println(match)
		if !match {
			return nil, errors.Parameter.WithMsg("产品id格式不对,格式为5到20位数字和英文字母组成的字符串")
		}
		find, err := l.CheckProductID(in)
		if err != nil {
			return nil, errors.System.AddDetail(err)
		} else if find {
			return nil, errors.Duplicate.WithMsgf("产品id重复:%s", in.ProductID).AddDetail("ProductID:" + in.ProductID)
		}
	} else {
		productID, err := relationDB.NewProductIDRepo(l.ctx).GenID(l.ctx)
		if err != nil {
			return nil, err
		}
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
	err = relationDB.NewProductInfoRepo(l.ctx).Insert(l.ctx, pi)
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, err
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, topics.DmProductInfoCreate, in.ProductID)
	if err != nil {
		l.Error(err)
	}
	return &dm.Empty{}, nil
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
	if err := l.svcCtx.StatusRepo.InitProduct(
		l.ctx, pi.ProductID); err != nil {
		l.Errorf("%s.StatusRepo.InitProduct failure,err:%v", utils.FuncName(), err)
		return errors.Database.AddDetail(err)
	}
	if err := l.svcCtx.SendRepo.InitProduct(
		l.ctx, pi.ProductID); err != nil {
		l.Errorf("%s.SendRepo.InitProduct failure,err:%v", utils.FuncName(), err)
		return errors.Database.AddDetail(err)
	}
	return nil
}

func CategorySchemaCreate(ctx context.Context, svcCtx *svc.ServiceContext, productID string, categoryID int64) error {
	//pc, err := relationDB.NewProductCategoryRepo(ctx).FindOne(ctx, categoryID)
	//if err != nil {
	//	return err
	//}
	//pcs, err := relationDB.NewProductCategorySchemaRepo(ctx).FindByFilter(ctx, relationDB.ProductCategorySchemaFilter{ProductCategoryIDs: utils.GetIDPath(pc.IDPath)}, nil)
	//if err != nil {
	//	return nil, err
	//}
	//ProductCategoryIDs = append(ProductCategoryIDs, utils.GetIDPath(pc.IDPath)...)
	return nil
}
