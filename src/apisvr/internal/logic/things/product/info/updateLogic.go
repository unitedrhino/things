package info

import (
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.ProductInfoUpdateReq) error {
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
	if req.Desc != nil {
		dmReq.Desc = &wrappers.StringValue{
			Value: *req.Desc,
		}
	}
	//if req.DevStatus != nil {
	//	dmReq.DevStatus = &wrappers.StringValue{
	//		Value: *req.DevStatus,
	//	}
	//}
	_, err := l.svcCtx.ProductM.ProductInfoUpdate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageProduct req=%v err=%+v", utils.FuncName(), utils.Fmt(req), er)
		return er
	}
	return nil
}
