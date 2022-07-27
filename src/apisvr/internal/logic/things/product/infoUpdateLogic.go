package product

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoUpdateLogic {
	return &InfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoUpdateLogic) InfoUpdate(req *types.ProductInfoUpdateReq) error {
	dmReq := &dm.ProductInfo{
		ProductID:    req.ProductID,
		ProductName:  req.ProductName,
		AuthMode:     req.AuthMode,
		DeviceType:   req.DeviceType,
		CategoryID:   req.CategoryID,
		NetType:      req.NetType,
		DataProto:    req.DataProto,
		AutoRegister: req.AutoRegister,
	}
	if req.Description != nil {
		dmReq.Description = &wrappers.StringValue{
			Value: *req.Description,
		}
	}
	//if req.DevStatus != nil {
	//	dmReq.DevStatus = &wrappers.StringValue{
	//		Value: *req.DevStatus,
	//	}
	//}
	_, err := l.svcCtx.DmRpc.ProductInfoUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageProduct|req=%v|err=%+v", utils.FuncName(), utils.Fmt(req), er)
		return er
	}
	return nil
}
