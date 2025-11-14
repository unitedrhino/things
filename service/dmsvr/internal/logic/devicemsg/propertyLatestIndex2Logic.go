package devicemsglogic

import (
	"context"
	"sync"

	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLatestIndex2Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyLatestIndex2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLatestIndex2Logic {
	return &PropertyLatestIndex2Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyLatestIndex2Logic) PropertyLatestIndex2(in *dm.PropertyLatestIndex2Req) (*dm.PropertyLogIndexResp, error) {
	var devMap = map[devices.Core][]string{}
	for _, dev := range in.Devices {
		d := devices.Core{
			ProductID:  dev.ProductID,
			DeviceName: dev.DeviceName,
		}
		_, ok := devMap[d]
		if !ok {
			devMap[d] = make([]string, 0)
		}
		devMap[d] = append(devMap[d], dev.DataID)
	}
	var ret []*dm.PropertyLogInfo
	var m sync.Mutex
	var wg sync.WaitGroup
	for dev, dataIDs := range devMap {
		wg.Add(1)
		utils.Go(l.ctx, func() {
			defer wg.Done()
			datas, err := NewPropertyLatestIndexLogic(l.ctx, l.svcCtx).PropertyLatestIndex(&dm.PropertyLatestIndexReq{
				ProductID:  dev.ProductID,
				DeviceName: dev.DeviceName,
				DataIDs:    dataIDs,
			})
			if err != nil {
				return
			}
			m.Lock()
			defer m.Unlock()
			for _, v := range datas.List {
				v.ProductID = dev.ProductID
				v.DeviceName = dev.DeviceName
				ret = append(ret, v)
			}
		})
	}
	wg.Wait()

	return &dm.PropertyLogIndexResp{List: ret}, nil
}
