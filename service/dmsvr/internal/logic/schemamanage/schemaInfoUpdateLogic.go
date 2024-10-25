package schemamanagelogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSchemaInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaInfoUpdateLogic {
	return &SchemaInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新设备物模型
func (l *SchemaInfoUpdateLogic) SchemaInfoUpdate(in *dm.SchemaInfo) (*dm.Empty, error) {
	// todo: add your logic here and delete this line

	return &dm.Empty{}, nil
}
