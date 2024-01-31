package schema

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TslImportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTslImportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TslImportLogic {
	return &TslImportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TslImportLogic) TslImport(req *types.ProductSchemaTslImportReq) error {
	dmReq := &dm.ProductSchemaTslImportReq{
		ProductID: req.ProductID, //产品id 只读
		Tsl:       req.Tsl,
	}
	_, err := l.svcCtx.ProductM.ProductSchemaTslImport(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ProductSchemaTslImport req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
