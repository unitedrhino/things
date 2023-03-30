package interact

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewActionReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionReadLogic {
	return &ActionReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActionReadLogic) ActionRead(req *types.DeviceInteractRespReadReq) (resp *types.DeviceInteractSendActionResp, err error) {
	dmReq := &di.RespReadReq{
		ProductID:   req.ProductID,
		DeviceName:  req.DeviceName,
		ClientToken: req.ClientToken,
	}
	dmResp, err := l.svcCtx.DeviceInteract.ActionRead(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ActionRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.DeviceInteractSendActionResp{
		ClientToken:  dmResp.ClientToken,
		OutputParams: dmResp.OutputParams,
		Status:       dmResp.Status,
		Code:         dmResp.Code,
	}, nil
}
