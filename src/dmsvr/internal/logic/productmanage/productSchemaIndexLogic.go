package productmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/repo/mysql"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProductSchemaIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProductSchemaIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductSchemaIndexLogic {
	return &ProductSchemaIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品信息列表
func (l *ProductSchemaIndexLogic) ProductSchemaIndex(in *dm.ProductSchemaIndexReq) (*dm.ProductSchemaIndexResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	schemas, err := l.svcCtx.ProductSchema.FindByFilter(l.ctx, mysql.ProductSchemaFilter{
		ProductID:   in.ProductID,
		Type:        in.Type,
		Tag:         in.Tag,
		Identifiers: in.Identifiers,
	}, def.PageInfo{Page: in.Page.GetPage(), Size: in.Page.GetSize()})
	if err != nil {
		if err == mysql.ErrNotFound {
			return nil, errors.NotFind
		}
		return nil, err
	}
	list := make([]*dm.ProductSchemaInfo, 0, len(schemas))
	for _, s := range schemas {
		list = append(list, ToProductSchemaRpc(s))
	}
	return &dm.ProductSchemaIndexResp{List: list}, nil
}
