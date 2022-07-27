package logic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/schema"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoDeleteLogic {
	return &ProductInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除设备
func (l *ProductInfoDeleteLogic) ProductInfoDelete(in *dm.ProductInfoDeleteReq) (*dm.Response, error) {
	pt, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	err = l.svcCtx.HubLogRepo.DropProduct(l.ctx, in.ProductID)
	if err != nil {
		l.Errorf("DelProduct|DeviceLogRepo|DropProduct|err=%+v", err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.DeviceDataRepo.DropProduct(l.ctx, pt, in.ProductID)
	if err != nil {
		l.Errorf("DelProduct|DeviceDataRepo|DropProduct|err=%+v", err)
		return nil, errors.Database.AddDetail(err)
	}
	l.svcCtx.SchemaRepo.ClearCache(l.ctx, in.ProductID)
	err = l.svcCtx.DmDB.Delete(l.ctx, in.ProductID)
	if err != nil {
		l.Errorf("DelProduct|Delete|err=%+v", err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.DataUpdate.TempModelUpdate(l.ctx, &schema.SchemaInfo{ProductID: in.ProductID})
	if err != nil {
		return nil, err
	}

	return &dm.Response{}, nil
}
