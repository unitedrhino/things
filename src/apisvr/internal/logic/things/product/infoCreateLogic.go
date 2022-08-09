package product

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoCreateLogic {
	return &InfoCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoCreateLogic) InfoCreate(req *types.ProductInfoCreateReq) error {
	dmReq := &dm.ProductInfo{
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
	_, err := l.svcCtx.DmRpc.ProductInfoCreate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.ManageProduct|req=%v|err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
