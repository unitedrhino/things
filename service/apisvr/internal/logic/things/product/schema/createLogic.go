package schema

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProductSchemaCreateReq) error {
	dmReq := &dm.ProductSchemaCreateReq{
		Info: ToSchemaInfoRpc(req.ProductSchemaInfo),
	}
	_, err := l.svcCtx.ProductM.ProductSchemaCreate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProductSchemaCreate req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
