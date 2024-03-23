package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
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

func (l *ReadLogic) Read(req *types.WithID) (resp *types.FirmwareInfo, err error) {
	var firmwareReadReq dm.WithID
	_ = utils.CopyE(&firmwareReadReq, &req)
	read, err := l.svcCtx.OtaM.OtaFirmwareInfoRead(l.ctx, &firmwareReadReq)
	logx.Infof("read:%+v", read)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareRead req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	var result = types.FirmwareInfo{FileList: []*types.FirmwareFile{}}
	_ = utils.CopyE(&result, &read)
	return &result, nil
}
