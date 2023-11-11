package task

import (
	"context"

	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

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

func (l *IndexLogic) Index(req *types.OtaTaskIndexReq) (resp *types.OtaTaskIndexResp, err error) {
	otaResp, err := l.svcCtx.OtaTaskM.OtaTaskIndex(l.ctx, &dm.OtaTaskIndexReq{
		FirmwareID: req.FirmwareID,
		Page:       logic.ToOtaPageRpc(req.Page),
	})
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.OtaTaskInfo, 0, len(otaResp.List))
	for _, v := range otaResp.List {
		pi := otaTaskInfoToApi(v)
		pis = append(pis, pi)
	}
	return &types.OtaTaskIndexResp{
		List:  pis,
		Total: otaResp.Total,
	}, nil
}
