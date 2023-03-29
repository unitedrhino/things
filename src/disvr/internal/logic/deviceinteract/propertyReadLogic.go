package deviceinteractlogic

import (
	"context"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyReadLogic {
	return &PropertyReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取异步调用设备属性的结果
func (l *PropertyReadLogic) PropertyRead(in *di.RespReadReq) (*di.SendPropertyResp, error) {
	resp, err := l.svcCtx.MsgThingRepo.GetResp(l.ctx, msgThing.TypeProperty,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, in.ClientToken)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	if resp == nil {
		return nil, errors.NotFind
	}
	return &di.SendPropertyResp{
		ClientToken: resp.ClientToken,
		Status:      resp.Status,
		Code:        resp.Code,
	}, nil
}
