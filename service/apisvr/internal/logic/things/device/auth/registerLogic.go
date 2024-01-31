package auth

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dgsvr/pb/dg"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.DeviceRegisterReq) (resp *types.DeviceRegisterResp, err error) {
	data, err := l.svcCtx.DeviceA.DeviceRegister(l.ctx, &dg.DeviceRegisterReq{
		ProductID:  req.ProductID,
		DeviceName: req.DeviceName,
		Nonce:      req.Nonce,
		Timestamp:  req.Timestamp,
		Signature:  req.Signature,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.DeviceRegister req=%v err=%+v", utils.FuncName(), req, er)
		return nil, err
	}
	return &types.DeviceRegisterResp{Len: data.Len, Payload: data.Payload}, nil
}
