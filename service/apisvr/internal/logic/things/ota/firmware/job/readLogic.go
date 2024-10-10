package job

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
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

func (l *ReadLogic) Read(req *types.WithID) (resp *types.OtaFirmwareJobInfo, err error) {
	var firmwareReadReq dm.WithID
	_ = utils.CopyE(&firmwareReadReq, &req)
	read, err := l.svcCtx.OtaM.OtaFirmwareJobRead(l.ctx, &firmwareReadReq)
	if err != nil {
		er := errors.Fmt(err)
		return nil, er
	}
	var result = types.OtaFirmwareJobInfo{}
	_ = utils.CopyE(&result, &read)
	utils.CopyE(&result.OtaFirmwareJobStatic, &read.Static)
	utils.CopyE(&result.OtaFirmwareJobDynamic, &read.Dynamic)
	return &result, nil
}
