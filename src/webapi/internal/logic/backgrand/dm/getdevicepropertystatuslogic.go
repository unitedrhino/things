package dm

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/dmsvr/dm"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDevicePropertyStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDevicePropertyStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetDevicePropertyStatusLogic {
	return GetDevicePropertyStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDevicePropertyStatusLogic) GetDevicePropertyStatus(req types.GetDevicePropertyStatusReq) (resp *types.GetDevicePropertyStatusResp, err error) {
	l.Infof("GetDevicePropertyStatus|req=%+v", req)
	tlReq := &dm.GetProductTemplateReq{
		ProductID: req.ProductID, //产品id
	}
	tlResp, err := l.svcCtx.DmRpc.GetProductTemplate(l.ctx, tlReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if tlResp.GetTemplate() == nil {
		return &types.GetDevicePropertyStatusResp{List: nil}, nil
	}
	template := tlResp.GetTemplate().GetValue()

	return
}
