package otaFirmware

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.FirmwareCreateReq) (resp *types.FirmwareResp, err error) {
	var firmwareCreateReq dm.OtaFirmwareCreateReq
	logx.Infof("req:%+v", req)
	_ = copier.Copy(&firmwareCreateReq, &req)
	firmwareCreateReq.FirmwareUdi = &wrappers.StringValue{
		Value: req.FirmwareUdi,
	}
	logx.Infof("firmwareCreateReq:%+v", &firmwareCreateReq)
	create, err := l.svcCtx.OtaFirmwareM.OtaFirmwareCreate(l.ctx, &firmwareCreateReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareCreate req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.FirmwareResp{FirmwareID: create.FirmwareID}, nil
}
