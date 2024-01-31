package info

import (
	"context"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/logic/things"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"sync"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

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

func (l *IndexLogic) Index(req *types.DeviceInfoIndexReq) (resp *types.DeviceInfoIndexResp, err error) {
	dmReq := &dm.DeviceInfoIndexReq{
		ProductID:         req.ProductID, //产品id
		AreaIDs:           req.AreaIDs,   //项目区域ids
		DeviceName:        req.DeviceName,
		Tags:              logic.ToTagsMap(req.Tags),
		Page:              logic.ToDmPageRpc(req.Page),
		Range:             req.Range,
		Position:          logic.ToDmPointRpc(req.Position),
		DeviceAlias:       req.DeviceAlias,
		IsOnline:          req.IsOnline,
		ProductCategoryID: req.ProductCategoryID,
	}
	dmResp, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceInfo req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.DeviceInfo, 0, len(dmResp.List))
	wait := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, v := range dmResp.List {
		wait.Add(1)
		info := v
		utils.Go(l.ctx, func() {
			defer wait.Done()
			pi := things.InfoToApi(l.ctx, l.svcCtx, info, req.WithProperties)
			mutex.Lock()
			defer mutex.Unlock()
			pis = append(pis, pi)
		})
	}
	wait.Wait()
	return &types.DeviceInfoIndexResp{
		Total: dmResp.Total,
		List:  pis,
		Num:   int64(len(pis)),
	}, nil
}
