package productmanagelogic

import (
	"context"
	schema2 "github.com/i-Things/things/shared/domain/schema"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/events"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/spf13/cast"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductSchemaUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaUpdateLogic {
	return &ProductSchemaUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProductSchemaUpdateLogic) ModifyProductSchema(in *dm.ProductSchemaUpdateReq, oldT *schema2.Model) (*dm.Response, error) {
	l.Infof("%s ProductID:%v", utils.FuncName(), in.Info.ProductID)
	newT, err := schema2.ValidateWithFmt([]byte(in.Info.Schema))
	if err != nil {
		return nil, err
	}
	err = schema2.CheckModify(oldT, newT)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.SchemaManaRepo.ModifyProduct(l.ctx, oldT, newT, in.Info.ProductID); err != nil {
		l.Errorf("%s ModifyProduct failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.SchemaRepo.Update(l.ctx, in.Info.ProductID, newT)
	if err != nil {
		l.Errorf("%s.Update err=%+v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	err = l.svcCtx.DataUpdate.ProductSchemaUpdate(l.ctx, &events.DataUpdateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}

func (l *ProductSchemaUpdateLogic) AddProductSchema(in *dm.ProductSchemaUpdateReq) (*dm.Response, error) {
	l.Infof("%s ProductID:%v", utils.FuncName(), in.Info.ProductID)
	_, err := l.svcCtx.ProductInfo.FindOne(l.ctx, in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.Parameter.AddDetail("not find ProductID id:" + cast.ToString(in.Info.ProductID))
		}
		return nil, errors.Database.AddDetail(err)
	}
	t, err := schema2.ValidateWithFmt([]byte(in.Info.Schema))
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.HubLogRepo.InitProduct(
		l.ctx, in.Info.ProductID); err != nil {
		l.Errorf("%s.DeviceLogRepo.InitProduct failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	if err := l.svcCtx.SchemaManaRepo.InitProduct(l.ctx, t, in.Info.ProductID); err != nil {
		l.Errorf("%s.SchemaManaRepo.InitProduct failure,err:%v", utils.FuncName(), err)
		return nil, errors.Database.AddDetail(err)
	}
	err = l.svcCtx.SchemaRepo.Insert(l.ctx, in.Info.ProductID, t)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.DataUpdate.ProductSchemaUpdate(l.ctx, &events.DataUpdateInfo{ProductID: in.Info.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, err
}

// 更新产品物模型
func (l *ProductSchemaUpdateLogic) ProductSchemaUpdate(in *dm.ProductSchemaUpdateReq) (*dm.Response, error) {
	pt, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.Info.ProductID)
	if err != nil {
		if err == mysql.ErrNotFound {
			return l.AddProductSchema(in)
		}
		return nil, errors.System.AddDetail(err)
	}

	return l.ModifyProductSchema(in, pt)
}
