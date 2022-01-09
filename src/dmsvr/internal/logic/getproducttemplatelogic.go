package logic

import (
	"context"

	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/go-things/things/src/dmsvr/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetProductTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductTemplateLogic {
	return &GetProductTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取产品信息
func (l *GetProductTemplateLogic) GetProductTemplate(in *dm.GetProductTemplateReq) (*dm.ProductTemplate, error) {
	// todo: add your logic here and delete this line

	return &dm.ProductTemplate{}, nil
}
