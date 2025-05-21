package msg

import (
	"context"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"sync"

	"gitee.com/unitedrhino/things/service/apisvr/internal/svc"
	"gitee.com/unitedrhino/things/service/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogMultiIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取单个id属性历史记录
func NewPropertyLogMultiIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogMultiIndexLogic {
	return &PropertyLogMultiIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PropertyLogMultiIndexLogic) PropertyLogMultiIndex(req *types.DeviceMsgPropertyLogMultiIndexReq) (resp *types.DeviceMsgPropertyMultiIndexResp, err error) {
	var wait sync.WaitGroup
	var ret = make([][]*types.DeviceMsgPropertyLogInfo, len(req.Reqs))
	var mutex sync.Mutex
	for i, r := range req.Reqs {
		wait.Add(1)
		ii := i
		rr := r
		utils.Go(l.ctx, func() {
			defer wait.Done()
			dmResp, err := l.svcCtx.DeviceMsg.PropertyLogIndex(l.ctx, utils.Copy[dm.PropertyLogIndexReq](rr))
			if err != nil {
				er := errors.Fmt(err)
				l.Errorf("%s.rpc.PropertyLogIndex req=%v err=%+v", utils.FuncName(), rr, er)
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			ret[ii] = utils.CopySlice[types.DeviceMsgPropertyLogInfo](dmResp.List)
		})
	}
	wait.Wait()
	return &types.DeviceMsgPropertyMultiIndexResp{
		Lists: ret,
	}, nil
}
