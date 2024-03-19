package job

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
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
	_ = copier.Copy(&firmwareReadReq, &req)
	read, err := l.svcCtx.OtaM.OtaFirmwareJobRead(l.ctx, &firmwareReadReq)
	if err != nil {
		er := errors.Fmt(err)
		return nil, er
	}
	var result = types.OtaFirmwareJobInfo{}
	_ = copier.Copy(&result, &read)
	copier.Copy(&result.OtaFirmwareJobStatic, &read.Static)
	copier.Copy(&result.OtaFirmwareJobDynamic, &read.Dynamic)
	return &result, nil
}
