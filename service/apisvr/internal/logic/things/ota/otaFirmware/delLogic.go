package otaFirmware

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelLogic {
	return &DelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelLogic) Del(req *types.FirmwareDeleteReq) error {
	var firmwareDeleteReq dm.OtaFirmwareDeleteReq
	_ = copier.Copy(&firmwareDeleteReq, &req)
	_, err := l.svcCtx.OtaFirmwareM.OtaFirmwareDelete(l.ctx, &firmwareDeleteReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareDelete req=%v err=%+v", utils.FuncName(), req, er)
		return er
	}
	return nil
}
