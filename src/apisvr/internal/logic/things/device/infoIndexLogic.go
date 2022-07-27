package device

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/dm"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type InfoIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InfoIndexLogic {
	return &InfoIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InfoIndexLogic) InfoIndex(req *types.DeviceInfoIndexReq) (resp *types.DeviceInfoIndexResp, err error) {
	dmReq := &dm.DeviceInfoIndexReq{
		ProductID: req.ProductID, //产品id
		Page: &dm.PageInfo{
			Page: req.Page.Page,
			Size: req.Page.Size,
		},
	}
	dmResp, err := l.svcCtx.DmRpc.DeviceInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s|rpc.GetDeviceInfo|req=%v|err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.DeviceInfo, 0, len(dmResp.List))
	for _, v := range dmResp.List {
		pi := deviceInfoToApi(v)
		pis = append(pis, pi)
	}
	return &types.DeviceInfoIndexResp{
		Total: dmResp.Total,
		List:  pis,
		Num:   int64(len(pis)),
	}, nil
}
