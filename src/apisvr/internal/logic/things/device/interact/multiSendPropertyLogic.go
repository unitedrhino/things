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

type MultiSendPropertyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiSendPropertyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiSendPropertyLogic {
	return &MultiSendPropertyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MultiSendPropertyLogic) MultiSendProperty(req *types.DeviceInteractMultiSendPropertyReq) (resp *types.DeviceInteractMultiSendPropertyResp, err error) {
	list := make([]*types.DeviceInteractMultiSendPropertyMsg, 0)

	dmReq := &di.MultiSendPropertyReq{
		ProductID:   req.ProductID,
		DeviceNames: req.DeviceNames,
		Data:        req.Data,
	}
	dmResp, err := l.svcCtx.DeviceInteract.MultiSendProperty(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.MultiSendProperty req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	if len(dmResp.List) > 0 {
		for _, v := range dmResp.List {
			list = append(list, &types.DeviceInteractMultiSendPropertyMsg{
				Code:        v.Code,
				Status:      v.Status,
				ClientToken: v.ClientToken,
				Data:        v.Data,
				ErrMsg:      v.ErrMsg,
				ErrCode:     v.ErrCode,
			})
		}
	}

	return &types.DeviceInteractMultiSendPropertyResp{List: list}, nil
}
