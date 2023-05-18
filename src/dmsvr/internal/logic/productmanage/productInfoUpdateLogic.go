package productmanagelogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/oss"
	"path"

	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoUpdateLogic {
	return &ProductInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductInfoUpdateLogic) setPoByPb(old *mysql.DmProductInfo, data *dm.ProductInfo) error {
	if data.Tags != nil {
		tags, err := json.Marshal(data.Tags)
		if err == nil {
			old.Tags = string(tags)
		}
	}
	if data.ProductName != "" {
		old.ProductName = data.ProductName
	}
	if data.AuthMode != def.Unknown {
		old.AuthMode = data.AuthMode
	}
	if data.Desc != nil {
		old.Desc = data.Desc.GetValue()
	}

	if data.AutoRegister != def.Unknown {
		old.AutoRegister = data.AutoRegister
	}
	if data.DevStatus != nil {
		old.DevStatus = data.DevStatus.GetValue()
	}

	if data.ProductName != "" {
		old.ProductName = data.ProductName
	}
	if data.ProductImg != "" {
		old.ProductImg = data.ProductImg
	}
	if data.ProductImg != "" && data.IsUpdateProductImg == true { //如果填了参数且不等于原来的,说明修改头像,需要处理
		si, err := oss.GetSceneInfo(data.ProductImg)
		if err != nil {
			return err
		}
		if !(si.Business == oss.BusinessProductManage && si.Scene == oss.SceneProductImg) {
			return errors.Parameter.WithMsg("产品图片的路径不对")
		}
		si.FilePath = data.ProductID + path.Ext(si.FilePath)
		nwePath, err := oss.GetFilePath(si, false)
		if err != nil {
			return err
		}
		path, err := l.svcCtx.OssClient.PublicBucket().CopyFromTempBucket(data.ProductImg, nwePath)
		if err != nil {
			return errors.System.AddDetail(err)
		}
		old.ProductImg = path
	}
	if data.AuthMode != 0 {
		old.AuthMode = data.AuthMode
	}
	if data.DeviceType != 0 {
		old.DeviceType = data.DeviceType
	}
	if data.CategoryID != 0 {
		old.CategoryID = data.CategoryID
	}
	if data.NetType != 0 {
		old.NetType = data.NetType
	}
	if data.DataProto != 0 {
		old.DataProto = data.DataProto
	}
	if data.AutoRegister != 0 {
		old.AutoRegister = data.AutoRegister
	}
	return nil
}

// 更新设备
func (l *ProductInfoUpdateLogic) ProductInfoUpdate(in *dm.ProductInfo) (*dm.Response, error) {
	po, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.ProductID))
		}
		return nil, errors.Database.AddDetail(err)
	}

	err = l.setPoByPb(po, in)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.ProductInfo.Update(l.ctx, po)
	if err != nil {
		l.Errorf("%s.Update err=%+v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}

	return &dm.Response{}, nil
}
