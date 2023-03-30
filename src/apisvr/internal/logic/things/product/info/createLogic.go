package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.ProductInfoCreateReq) error {
	dmReq := &dm.ProductInfo{
		ProductName:  req.ProductName,
		AuthMode:     req.AuthMode,
		DeviceType:   req.DeviceType,
		CategoryID:   req.CategoryID,
		NetType:      req.NetType,
		DataProto:    req.DataProto,
		AutoRegister: req.AutoRegister,
		Desc:         utils.ToRpcNullString(req.Desc),
		Tags:         logic.ToTagsMap(req.Tags),
	}
	//if req.DevStatus != nil {
	//	dmReq.DevStatus = &wrappers.StringValue{
	//		Value: *req.DevStatus,
	//	}
	//}
	_, err := l.svcCtx.ProductM.ProductInfoCreate(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ManageProduct req=%v err=%v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
