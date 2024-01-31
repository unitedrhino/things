package firmware

import (
	"context"

	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic"
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

func (l *IndexLogic) Index(req *types.OtaFirmwareIndexReq) (resp *types.OtaFirmwareIndexResp, err error) {
	indexReq := &dm.FirmwareInfoIndexReq{
		ProductID: req.ProductID,
		Page:      logic.ToOtaPageRpc(req.Page),
	}
	dmResp, err := l.svcCtx.FirmwareM.FirmwareInfoIndex(l.ctx, indexReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.OtaFirmwareIndex, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		pi := firmwareInfoToApi(v)
		pis = append(pis, pi)
	}
	return &types.OtaFirmwareIndexResp{
		Total: dmResp.Total,
		List:  pis,
	}, nil
}
