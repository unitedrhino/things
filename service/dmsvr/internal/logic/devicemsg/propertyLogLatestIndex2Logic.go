package devicemsglogic

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/share/devices"
	"sync"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogLatestIndex2Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyLogLatestIndex2Logic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogLatestIndex2Logic {
	return &PropertyLogLatestIndex2Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyLogLatestIndex2Logic) PropertyLogLatestIndex2(in *dm.PropertyLogLatestIndex2Req) (*dm.PropertyLogIndexResp, error) {
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
			datas, err := NewPropertyLogLatestIndexLogic(l.ctx, l.svcCtx).PropertyLogLatestIndex(&dm.PropertyLogLatestIndexReq{
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
