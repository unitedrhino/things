package otaFirmware

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.FirmwareReadReq) (resp *types.FirmwareReadResp, err error) {
	var firmwareReadReq dm.OtaFirmwareReadReq
	_ = copier.Copy(&firmwareReadReq, &req)
	read, err := l.svcCtx.OtaFirmwareM.OtaFirmwareRead(l.ctx, &firmwareReadReq)
	logx.Infof("read:%+v", read)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	var result types.FirmwareReadResp
	_ = copier.Copy(&result, &read)
	logx.Infof("resp:%+v", result)
	return &result, nil
}
