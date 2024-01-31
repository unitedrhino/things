package productmanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/shared/def"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/oss"
	"gitee.com/i-Things/core/shared/oss/common"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
}

func NewProductInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoUpdateLogic {
	return &ProductInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
	}
}

func (l *ProductInfoUpdateLogic) setPoByPb(old *relationDB.DmProductInfo, data *dm.ProductInfo) error {
	if data.Tags != nil {
		old.Tags = data.Tags
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
	if data.ProductImg != "" && data.IsUpdateProductImg == true { //如果填了参数且不等于原来的,说明修改头像,需要处理
		if old.ProductImg != "" {
			err := l.svcCtx.OssClient.PrivateBucket().Delete(l.ctx, old.ProductImg, common.OptionKv{})
			if err != nil {
				l.Errorf("Delete file err path:%v,err:%v", old.ProductImg, err)
			}
		}
		nwePath := oss.GenFilePath(l.ctx, l.svcCtx.Config.Name, oss.BusinessProductManage, oss.SceneProductImg, fmt.Sprintf("%s/%s", data.ProductID, oss.GetFileNameWithPath(data.ProductImg)))
		path, err := l.svcCtx.OssClient.PrivateBucket().CopyFromTempBucket(data.ProductImg, nwePath)
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
	po, err := l.PiDB.FindOneByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddDetail("not find Product_id id:" + cast.ToString(in.ProductID))
		}
		return nil, err
	}

	err = l.setPoByPb(po, in)
	if err != nil {
		return nil, err
	}
	err = l.PiDB.Update(l.ctx, po)
	if err != nil {
		l.Errorf("%s.Update err=%+v", utils.FuncName(), err)
		if errors.Cmp(err, errors.Duplicate) {
			return nil, errors.Duplicate.WithMsgf("产品名称重复:%s", in.ProductName)
		}
		return nil, err
	}

	return &dm.Response{}, nil
}
