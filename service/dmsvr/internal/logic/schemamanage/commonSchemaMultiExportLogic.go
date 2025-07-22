package schemamanagelogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommonSchemaMultiExportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCommonSchemaMultiExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommonSchemaMultiExportLogic {
	return &CommonSchemaMultiExportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CommonSchemaMultiExportLogic) CommonSchemaMultiExport(in *dm.CommonSchemaExportReq) (*dm.CommonSchemaExportResp, error) {
	pos, err := relationDB.NewCommonSchemaRepo(l.ctx).FindByFilter(l.ctx, relationDB.CommonSchemaFilter{Identifiers: in.Identifiers}, nil)
	if err != nil {
		return nil, err
	}
	list := make([]*dm.CommonSchemaInfo, 0, len(pos))
	for _, s := range pos {
		list = append(list, ToCommonSchemaRpc(s))
	}
	return &dm.CommonSchemaExportResp{Schemas: utils.MarshalNoErr(list)}, nil
}
