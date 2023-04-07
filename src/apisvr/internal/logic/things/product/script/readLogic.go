package script

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.ProductScriptReadReq) (resp *types.ProductScript, err error) {
	dmResp, err := l.svcCtx.ProductM.ProductScriptRead(l.ctx, &dm.ProductScriptReadReq{ProductID: req.ProductID})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s rpc.ProductScriptRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.ProductScript{
		ProductID: dmResp.ProductID,
		Script:    dmResp.Script,
		Lang:      dmResp.Lang,
	}, nil
}
