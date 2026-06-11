package pair

import (
	"context"
	"strconv"

	"gitee.com/unitedrhino/share/ctxs"
	shareerr "gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmLogic {
	return &ConfirmLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmLogic) Confirm(req *types.DevicePairConfirmReq) (*types.DevicePairConfirmResp, error) {
	if req.ProductID == "" || req.Mac == "" || req.GrantToken == "" || req.PairAckPayload == "" {
		return nil, shareerr.Parameter.WithMsg("productID、mac、grantToken和pairAckPayload必填")
	}
	product, err := l.svcCtx.ProductM.ProductInfoRead(l.ctx, &dm.ProductInfoReadReq{
		ProductID: req.ProductID,
	})
	if err != nil {
		return nil, err
	}
	if _, err := DecodeProductMK(product.Secret); err != nil {
		return nil, shareerr.Parameter.WithMsg("产品密钥必须是32位hex MK")
	}
	mac, _, err := NormalizeMAC(req.Mac)
	if err != nil {
		return nil, shareerr.Parameter.WithMsg("mac格式不正确")
	}
	deviceName := req.DeviceName
	if deviceName == "" {
		deviceName = mac
	}
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	grant, err := VerifyGrant(VerifyGrantInput{
		Token:      req.GrantToken,
		SigningKey: GrantSigningKey(req.ProductID, product.Secret),
		ProductID:  req.ProductID,
		MAC:        mac,
		DeviceName: deviceName,
		UserID:     strconv.FormatInt(uc.UserID, 10),
	})
	if err != nil {
		return nil, pairError(err)
	}
	ack, err := VerifyPairAck(req.PairAckPayload, mac, grant.PairKeyHex, grant.ObservedBindEpoch)
	if err != nil {
		return nil, pairError(err)
	}
	_, err = l.svcCtx.DeviceM.DeviceInfoBind(l.ctx, &dm.DeviceInfoBindReq{
		Device: &dm.DeviceCore{
			ProductID:  req.ProductID,
			DeviceName: deviceName,
		},
		AreaID:          req.AreaID,
		IsIgnoreOffline: true,
	})
	if err != nil {
		return nil, err
	}
	return &types.DevicePairConfirmResp{
		ProductID:  req.ProductID,
		Mac:        mac,
		DeviceName: deviceName,
		BindEpoch:  ack.BindEpoch,
		BleSecVer:  2,
		BlePairKey: grant.PairKeyHex,
		Message:    "bind_confirmed",
	}, nil
}

func pairError(err error) error {
	switch err {
	case ErrGrantTokenExpired:
		return shareerr.Permissions.WithMsg("grant_token_expired")
	case ErrGrantTokenMismatch:
		return shareerr.Permissions.WithMsg("grant_token_mismatch")
	case ErrPairAckAuthInvalid:
		return shareerr.Permissions.WithMsg("pair_ack_auth_invalid")
	case ErrBindEpochTooOld:
		return shareerr.Permissions.WithMsg("bind_epoch_too_old")
	case ErrInvalidGrantToken:
		return shareerr.Permissions.WithMsg("invalid_grant_token")
	default:
		return shareerr.Parameter.WithMsg(err.Error())
	}
}
