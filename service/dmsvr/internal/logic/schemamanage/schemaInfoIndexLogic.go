package schemamanagelogic

import (
	"context"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSchemaInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaInfoIndexLogic {
	return &SchemaInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备物模型列表
func (l *SchemaInfoIndexLogic) SchemaInfoIndex(in *dm.SchemaInfoIndexReq) (*dm.SchemaInfoIndexResp, error) {
	// todo: add your logic here and delete this line

	return &dm.SchemaInfoIndexResp{}, nil
}
