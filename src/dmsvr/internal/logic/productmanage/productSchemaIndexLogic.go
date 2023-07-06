package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PsDb *relationDB.ProductSchemaRepo
}

func NewProductSchemaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaIndexLogic {
	return &ProductSchemaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PsDb:   relationDB.NewProductSchemaRepo(ctx),
	}
}

// 获取产品信息列表
func (l *ProductSchemaIndexLogic) ProductSchemaIndex(in *dm.ProductSchemaIndexReq) (*dm.ProductSchemaIndexResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	filter := relationDB.ProductSchemaFilter{
		ProductID:   in.ProductID,
		Type:        in.Type,
		Tag:         in.Tag,
		Identifiers: in.Identifiers,
	}
	schemas, err := l.PsDb.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind
		}
		return nil, err
	}
	total, err := l.PsDb.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	list := make([]*dm.ProductSchemaInfo, 0, len(schemas))
	for _, s := range schemas {
		list = append(list, ToProductSchemaRpc(s))
	}
	return &dm.ProductSchemaIndexResp{List: list, Total: total}, nil
}
