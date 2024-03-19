package info

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"
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

func (l *CreateLogic) Create(req *types.FirmwareCreateReq) (resp *types.WithID, err error) {
	var firmwareCreateReq dm.OtaFirmwareInfoCreateReq
	logx.Infof("req:%+v", req)
	_ = copier.Copy(&firmwareCreateReq, &req)
	firmwareCreateReq.Extra = utils.ToRpcNullString(req.Extra)
	logx.Infof("firmwareCreateReq:%+v", &firmwareCreateReq)
	create, err := l.svcCtx.OtaM.OtaFirmwareInfoCreate(l.ctx, &firmwareCreateReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareCreate req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.WithID{ID: create.Id}, nil
}
