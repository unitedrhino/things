package otaFirmware

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

func (l *IndexLogic) Index(req *types.FirmwareIndexReq) (resp *types.FirmwareIndexResp, err error) {
	var firmwareIndexReq dm.OtaFirmwareIndexReq
	_ = copier.Copy(&firmwareIndexReq, &req)
	index, err := l.svcCtx.OtaFirmwareM.OtaFirmwareIndex(l.ctx, &firmwareIndexReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.OtaFirmwareIndex req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	var list []types.FirmwareInfo
	_ = copier.Copy(&list, &index.List)
	return &types.FirmwareIndexResp{List: list, Total: index.Total}, nil
}
