package job

import (
	"context"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/jinzhu/copier"

	"github.com/i-Things/things/service/apisvr/internal/svc"
	"github.com/i-Things/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyLogic {
	return &VerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyLogic) Verify(req *types.OtaFirmwareVerifyReq) (resp *types.UpgradeJobResp, err error) {
	var otaFirmwareVerifyReq dm.OtaFirmwareVerifyReq
	_ = copier.Copy(&otaFirmwareVerifyReq, &req)
	create, err := l.svcCtx.OtaJobM.OtaVerifyJobCreate(l.ctx, &otaFirmwareVerifyReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaVerifyJobCreate req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.UpgradeJobResp{JobID: create.JobId, UtcCreate: create.UtcCreate}, nil
}
