package auth

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
	res, err := l.svcCtx.DeviceA.DeviceRegister(l.ctx, &dm.DeviceRegisterReq{
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
	return &types.DeviceRegisterResp{Psk: res.Psk}, nil
}
