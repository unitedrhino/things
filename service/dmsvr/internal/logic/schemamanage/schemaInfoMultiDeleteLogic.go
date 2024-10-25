package schemamanagelogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaInfoMultiDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSchemaInfoMultiDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaInfoMultiDeleteLogic {
	return &SchemaInfoMultiDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除设备物模型
func (l *SchemaInfoMultiDeleteLogic) SchemaInfoMultiDelete(in *dm.SchemaInfoMultiDeleteReq) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
