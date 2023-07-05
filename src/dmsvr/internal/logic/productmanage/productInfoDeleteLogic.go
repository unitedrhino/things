package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/events/topics"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDb *relationDB.ProductInfoRepo
}

func NewProductInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoDeleteLogic {
	return &ProductInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDb:   relationDB.NewProductInfoRepo(ctx),
	}
}

// 删除设备
func (l *ProductInfoDeleteLogic) ProductInfoDelete(in *dm.ProductInfoDeleteReq) (*dm.Response, error) {
	err := l.Check(in)
	if err != nil {
		return nil, err
	}
	err = l.DropProduct(in)
	if err != nil {
		return nil, err
	}
	err = l.PiDb.DeleteByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	err = l.svcCtx.DataUpdate.ProductSchemaUpdate(l.ctx, &events.DeviceUpdateInfo{ProductID: in.ProductID})
	if err != nil {
		return nil, err
	}

	return &dm.Response{}, nil
}
func (l *ProductInfoDeleteLogic) DropProduct(in *dm.ProductInfoDeleteReq) error {
	pt, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	err = l.svcCtx.HubLogRepo.DropProduct(l.ctx, in.ProductID)
	if err != nil {
		l.Errorf("%s.HubLogRepo.DeleteProduct err=%v", utils.FuncName(), utils.Fmt(err))
		return err
	}
	err = l.svcCtx.SDKLogRepo.DropProduct(l.ctx, in.ProductID)
	if err != nil {
		l.Errorf("%s.SDKLogRepo.DeleteProduct err=%v", utils.FuncName(), utils.Fmt(err))
		return err
	}
	err = l.svcCtx.SchemaManaRepo.DeleteProduct(l.ctx, pt, in.ProductID)
	if err != nil {
		l.Errorf("%s.SchemaManaRepo.DeleteProduct err=%v", utils.FuncName(), utils.Fmt(err))
		return err
	}
	//todo 需要删除物模型的数据
	err = l.svcCtx.SchemaRepo.ClearCache(l.ctx, in.ProductID)
	if err != nil {
		l.Errorf("%s.SchemaRepo.ClearCache err=%v", utils.FuncName(), utils.Fmt(err))
		return err
	}
	l.svcCtx.Bus.Publish(l.ctx, topics.DmProductInfoDelete, in.ProductID)
	return nil
}
func (l *ProductInfoDeleteLogic) Check(in *dm.ProductInfoDeleteReq) error {
	count, err := l.svcCtx.DeviceInfo.CountByFilter(l.ctx, mysql.DeviceFilter{ProductID: in.ProductID})
	if err != nil {
		l.Errorf("%s.CountByFilter err=%v", utils.FuncName(), utils.Fmt(err))
		return errors.Database.AddDetail(err)
	}
	if count > 0 {
		return errors.NotEmpty.WithMsg("该产品下还有设备,不可删除")
	}
	return nil
}
