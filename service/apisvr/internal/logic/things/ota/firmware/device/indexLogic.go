package device

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.OtaFirmwareDeviceIndexReq) (resp *types.OtaFirmwareDeviceIndexResp, err error) {
	index, err := l.svcCtx.OtaM.OtaFirmwareDeviceIndex(l.ctx, utils.Copy[dm.OtaFirmwareDeviceIndexReq](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return utils.Copy[types.OtaFirmwareDeviceIndexResp](index), nil
}
