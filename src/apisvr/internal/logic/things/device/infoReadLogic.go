package device

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoReadLogic {
	return &InfoReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoReadLogic) InfoRead(req *types.DeviceInfoReadReq) (resp *types.DeviceInfo, err error) {
	dmResp, err := l.svcCtx.DmRpc.DeviceInfoRead(l.ctx,
		&dm.DeviceInfoReadReq{ProductID: req.ProductID, DeviceName: req.DeviceName})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return deviceInfoToApi(dmResp), nil
}
