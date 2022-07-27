package logic

import (
	"context"
	"database/sql"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"
	"time"

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

func (l *ProductInfoUpdateLogic) UpdateProductInfo(old *mysql.ProductInfo, data *dm.ProductInfo) {
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

// 更新设备
func (l *ProductInfoUpdateLogic) ProductInfoUpdate(in *dm.ProductInfo) (*dm.Response, error) {
	pi, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.ProductID))
		}
		return nil, errors.Database.AddDetail(err)
	}
	l.UpdateProductInfo(pi, in)

	err = l.svcCtx.ProductInfo.Update(l.ctx, pi)
	if err != nil {
		l.Errorf("ModifyProduct|ProductInfo|Update|err=%+v", err)
		return nil, errors.Database.AddDetail(err)
	}

	return &dm.Response{}, nil
}
