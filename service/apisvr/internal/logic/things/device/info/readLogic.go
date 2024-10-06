package info

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/apisvr/internal/logic/things"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"gitee.com/i-Things/things/service/apisvr/internal/svc"
	"gitee.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.DeviceInfoReadReq) (resp *types.DeviceInfo, err error) {
	dmResp, err := l.svcCtx.DeviceM.DeviceInfoRead(l.ctx,
		&dm.DeviceInfoReadReq{ProductID: req.ProductID, DeviceName: req.DeviceName, WithGateway: req.WithGateway})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceInfo req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return things.InfoToApi(l.ctx, l.svcCtx, dmResp, things.DeviceInfoWith{Owner: req.WithOwner, Properties: req.WithProperties, Profiles: req.WithProfiles}), nil
}
