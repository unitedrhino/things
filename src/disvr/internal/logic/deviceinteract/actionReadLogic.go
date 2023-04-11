package deviceinteractlogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/src/disvr/internal/repo/cache"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActionReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewActionReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActionReadLogic {
	return &ActionReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取异步调用设备行为的结果
func (l *ActionReadLogic) ActionRead(in *di.RespReadReq) (*di.SendActionResp, error) {
	resp, err := cache.GetDeviceMsg[msgThing.Resp](l.ctx, l.svcCtx.Store, deviceMsg.RespMsg, devices.Thing, msgThing.TypeAction,
		devices.Core{ProductID: in.ProductID, DeviceName: in.DeviceName}, in.ClientToken)
	if err != nil {
		return nil, errors.Database.AddDetail(err)
	}
	if resp == nil {
		return nil, errors.NotFind
	}
	respParam, err := json.Marshal(resp.Response)
	if err != nil {
		return nil, errors.RespParam.AddDetailf("SendAction get device resp not right:%+v", resp.Response)
	}
	return &di.SendActionResp{
		ClientToken:  resp.ClientToken,
		Status:       resp.Status,
		Code:         resp.Code,
		OutputParams: string(respParam),
	}, nil
}
