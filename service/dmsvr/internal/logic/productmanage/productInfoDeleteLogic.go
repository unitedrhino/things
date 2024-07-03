package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	DiDB *relationDB.DeviceInfoRepo
}

func NewProductInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductInfoDeleteLogic {
	return &ProductInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

// 删除设备
func (l *ProductInfoDeleteLogic) ProductInfoDelete(in *dm.ProductInfoDeleteReq) (*dm.Empty, error) {
	err := l.Check(in)
	if err != nil {
		return nil, err
	}
	err = relationDB.NewProductSchemaRepo(l.ctx).DeleteByFilter(l.ctx, relationDB.ProductSchemaFilter{ProductID: in.ProductID})
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	err = l.DropProduct(in)
	if err != nil {
		return nil, err
	}
	err = l.PiDB.DeleteByFilter(l.ctx, relationDB.ProductFilter{ProductIDs: []string{in.ProductID}})
	if err != nil {
		l.Errorf("%s.Delete err=%v", utils.FuncName(), utils.Fmt(err))
		return nil, err
	}
	return &dm.Empty{}, nil
}
func (l *ProductInfoDeleteLogic) DropProduct(in *dm.ProductInfoDeleteReq) error {
	pt, err := l.svcCtx.SchemaRepo.GetData(l.ctx, in.ProductID)
	if err != nil {
		return errors.System.AddDetail(err)
	}
	err = l.svcCtx.HubLogRepo.DeleteProduct(l.ctx, in.ProductID)
	if err != nil {
		l.Errorf("%s.HubLogRepo.DeleteProduct err=%v", utils.FuncName(), utils.Fmt(err))
		return err
	}
	err = l.svcCtx.SDKLogRepo.DeleteProduct(l.ctx, in.ProductID)
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
	err = l.svcCtx.SchemaRepo.SetData(l.ctx, in.ProductID, nil)
	if err != nil {
		l.Errorf("%s.SchemaRepo.ClearCache err=%v", utils.FuncName(), utils.Fmt(err))
		return err
	}
	err = l.svcCtx.ProductCache.SetData(l.ctx, in.ProductID, nil)
	if err != nil {
		l.Error(err)
	}
	err = l.svcCtx.FastEvent.Publish(l.ctx, eventBus.DmProductInfoDelete, in.ProductID)
	if err != nil {
		l.Error(err)
	}
	return nil
}
func (l *ProductInfoDeleteLogic) Check(in *dm.ProductInfoDeleteReq) error {
	count, err := l.DiDB.CountByFilter(l.ctx, relationDB.DeviceFilter{ProductID: in.ProductID})
	if err != nil {
		l.Errorf("%s.CountByFilter err=%v", utils.FuncName(), utils.Fmt(err))
		return err
	}
	if count > 0 {
		return errors.NotEmpty.WithMsg("该产品下还有设备,不可删除")
	}
	return nil
}
