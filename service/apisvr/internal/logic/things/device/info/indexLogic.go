package info

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/apisvr/internal/logic/things"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"sync"

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
		ctx:    ctxs.WithDefaultRoot(ctx),
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.DeviceInfoIndexReq) (resp *types.DeviceInfoIndexResp, err error) {
	dmResp, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, utils.Copy[dm.DeviceInfoIndexReq](req))
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.GetDeviceInfo req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	pis := make([]*types.DeviceInfo, 0, len(dmResp.List))
	var piMap = map[int64]*types.DeviceInfo{}
	wait := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, v := range dmResp.List {
		wait.Add(1)
		info := v
		utils.Go(l.ctx, func() {
			defer wait.Done()
			pi := things.InfoToApi(l.ctx, l.svcCtx, info, things.DeviceInfoWith{WithArea: req.WithArea, WithOwner: req.WithOwner, WithProperties: req.WithProperties, WithProfiles: req.WithProfiles})
			mutex.Lock()
			defer mutex.Unlock()
			piMap[pi.ID] = pi
		})
	}
	wait.Wait()
	for _, v := range dmResp.List {
		pis = append(pis, piMap[v.Id])
	}
	return &types.DeviceInfoIndexResp{
		Total: dmResp.Total,
		List:  pis,
	}, nil
}
