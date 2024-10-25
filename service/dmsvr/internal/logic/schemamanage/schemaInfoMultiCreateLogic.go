package schemamanagelogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaInfoMultiCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSchemaInfoMultiCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaInfoMultiCreateLogic {
	return &SchemaInfoMultiCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量新增物模型,只新增没有的,已有的不处理
func (l *SchemaInfoMultiCreateLogic) SchemaInfoMultiCreate(in *dm.SchemaInfoMultiCreateReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
