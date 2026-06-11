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

type GrantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGrantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GrantLogic {
	return &GrantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GrantLogic) Grant(req *types.DevicePairGrantReq) (*types.DevicePairGrantResp, error) {
	if req.ProductID == "" || req.Mac == "" {
		return nil, shareerr.Parameter.WithMsg("productID和mac必填")
	}
	product, err := l.svcCtx.ProductM.ProductInfoRead(l.ctx, &dm.ProductInfoReadReq{
		ProductID: req.ProductID,
	})
	if err != nil {
		return nil, err
	}
	mk, err := DecodeProductMK(product.Secret)
	if err != nil {
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
	grant, err := BuildGrant(GrantInput{
		ProductID:         req.ProductID,
		MAC:               mac,
		DeviceName:        deviceName,
		UserID:            strconv.FormatInt(uc.UserID, 10),
		ObservedBindEpoch: req.ObservedBindEpoch,
		MK:                mk,
		SigningKey:        GrantSigningKey(req.ProductID, product.Secret),
	})
	if err != nil {
		return nil, shareerr.Parameter.WithMsg(err.Error())
	}
	return &types.DevicePairGrantResp{
		ProductID:  grant.ProductID,
		Mac:        grant.MAC,
		DeviceName: grant.DeviceName,
		GrantToken: grant.GrantToken,
		Nonce:      grant.Nonce,
		AuthTag:    grant.AuthTag,
		TtlSec:     grant.TTLSec,
	}, nil
}
