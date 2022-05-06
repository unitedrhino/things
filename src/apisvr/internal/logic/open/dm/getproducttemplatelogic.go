package dm

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetProductTemplateLogic {
	return GetProductTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductTemplateLogic) GetProductTemplate(req types.GetProductTemplateReq) (resp *types.ProductTemplate, err error) {
	l.Infof("GetProductTemplate|req=%+v", req)
	dmReq := &dm.GetProductTemplateReq{
		ProductID: req.ProductID, //产品id
	}
	dmResp, err := l.svcCtx.DmRpc.GetProductTemplate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	resp = types.ProductTemplateToApi(dmResp)
	return resp, nil
}
