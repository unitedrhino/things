package msg

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShadowIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShadowIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShadowIndexLogic {
	return &ShadowIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ShadowIndexLogic) ShadowIndex(req *types.DeviceMsgPropertyLogLatestIndexReq) (resp *types.DeviceMsgShadowIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceMsg.ShadowIndex(l.ctx, &dm.PropertyLogLatestIndexReq{
		DeviceName: req.DeviceName,
		ProductID:  req.ProductID,
		DataIDs:    req.DataIDs,
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.ShadowIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	info := make([]*types.DeviceMsgShadowIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		info = append(info, &types.DeviceMsgShadowIndex{
			UpdatedDeviceTime: v.UpdatedDeviceTime,
			DataID:            v.DataID,
			Value:             v.Value,
		})
	}
	return &types.DeviceMsgShadowIndexResp{
		List: info,
	}, nil
}
