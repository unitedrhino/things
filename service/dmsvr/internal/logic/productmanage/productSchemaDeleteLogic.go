package productmanagelogic

import (
	"context"
	"gitee.com/i-Things/share/domain/schema"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/events"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PsDB *relationDB.ProductSchemaRepo
}

func NewProductSchemaDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaDeleteLogic {
	return &ProductSchemaDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PsDB:   relationDB.NewProductSchemaRepo(ctx),
	}
}

// 删除产品
func (l *ProductSchemaDeleteLogic) ProductSchemaDelete(in *dm.ProductSchemaDeleteReq) (*dm.Response, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	po, err := l.PsDB.FindOneByFilter(l.ctx, relationDB.ProductSchemaFilter{
		ProductID: in.ProductID, Identifiers: []string{in.Identifier},
	})
	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.Parameter.AddMsg("标识符未找到")
		}
		return nil, err
	}
	if schema.AffordanceType(po.Type) == schema.AffordanceTypeProperty {
		if err := l.svcCtx.SchemaManaRepo.DeleteProperty(l.ctx, in.ProductID, in.Identifier); err != nil {
			l.Errorf("%s.DeleteProperty failure,err:%v", utils.FuncName(), err)
			return nil, errors.Database.AddDetail(err)
		}
	}
	err = l.PsDB.DeleteByFilter(l.ctx, relationDB.ProductSchemaFilter{ID: po.ID})
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.ServerMsg.Publish(l.ctx, eventBus.DmProductSchemaUpdate, &events.DeviceUpdateInfo{ProductID: in.ProductID})
	if err != nil {
		return nil, err
	}
	return &dm.Response{}, nil
}
