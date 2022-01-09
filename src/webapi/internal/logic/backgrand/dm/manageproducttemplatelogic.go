package dm

import (
	"context"
	"github.com/go-things/things/shared/errors"
	"github.com/go-things/things/shared/utils"
	"github.com/go-things/things/src/dmsvr/dm"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/go-things/things/src/webapi/internal/svc"
	"github.com/go-things/things/src/webapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ManageProductTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewManageProductTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) ManageProductTemplateLogic {
	return ManageProductTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ManageProductTemplateLogic) ManageProductTemplate(req types.ManageProductTemplateReq) (resp *types.ProductTemplate, err error) {
	l.Infof("ManageProduct|req=%+v", req)
	dmReq := &dm.ManageProductTemplateReq{
		Info: &dm.ProductTemplate{
			ProductID: req.Info.ProductID, //产品id 只读
		},
	}
	if req.Info.Template != nil {
		dmReq.Info.Template = &wrappers.StringValue{
			Value: *req.Info.Template,
		}
	}
	dmResp, err := l.svcCtx.DmRpc.ManageProductTemplate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageProductTemplate|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return types.ProductTemplateToApi(dmResp), nil
}
