package devicemsglogic

import (
	"context"
	"sync"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"
	"golang.org/x/sync/errgroup"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyAgg2IndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyAgg2IndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyAgg2IndexLogic {
	return &PropertyAgg2IndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PropertyAgg2IndexLogic) PropertyAgg2Index(in *dm.PropertyAgg2IndexReq) (*dm.PropertyAggIndexResp, error) {
	var (
		//diDatas    []*dm.PropertyLogInfo
		dd  = l.svcCtx.SchemaManaRepo
		err error
	)
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if !uc.IsAdmin {
		for _, dev := range in.Aggs {
			_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
				ProductID:  dev.Device.ProductID,
				DeviceName: dev.Device.DeviceName,
			}, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	var diDatas []*dm.PropertyAggResp
	var mutex sync.Mutex
	var wg errgroup.Group
	for _, agg2 := range in.Aggs {
		agg := agg2
		wg.Go(func() error {
			defer utils.Recover(l.ctx)
			t, err := l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, devices.Core{ProductID: agg.Device.ProductID, DeviceName: agg.Device.DeviceName})
			if err != nil {
				return err
			}
			dds, err := dd.GetPropertyAgg(l.ctx, t, msgThing.FilterAggOpt{
				TimeStart: in.TimeStart,
				TimeEnd:   in.TimeEnd,
				Filter: msgThing.Filter{
					ProductID:    agg.Device.ProductID,
					DeviceNames:  []string{agg.Device.DeviceName},
					Interval:     in.Interval,
					IntervalUnit: def.TimeUnit(in.IntervalUnit),
					PartitionBy:  in.PartitionBy,
				}, Aggs: []msgThing.PropertyAgg{utils.Copy2[msgThing.PropertyAgg](agg)},
			})
			if err != nil {
				l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
				return err
			}
			mutex.Lock()
			defer mutex.Unlock()
			for _, devData := range dds {
				diData := dm.PropertyAggResp{
					ProductID:  agg.Device.ProductID,
					DeviceName: agg.Device.DeviceName,
				}
				for _, v1 := range devData.Values {
					var dv = dm.PropertyAggRespDetail{
						DataID:     v1.Identifier,
						DataName:   schema.GetDataName(t, v1.Identifier),
						TimeWindow: v1.TsWindow.UnixMilli(),
						Values:     map[string]*dm.PropertyAggRespDataDetail{}}
					for k2, v2 := range v1.Values {
						dv2 := dm.PropertyAggRespDataDetail{Timestamp: v2.TimeStamp.UnixMilli()}
						if dv2.Timestamp < 0 {
							dv2.Timestamp = 0
						}
						var payload []byte
						if param, ok := v2.Param.(string); ok {
							payload = []byte(param)
						} else {
							payload = []byte(utils.ToString(v2.Param))
						}
						dv2.Value = string(payload)
						dv.Values[k2] = &dv2
					}
					diData.Values = append(diData.Values, &dv)
				}
				diDatas = append(diDatas, &diData)
			}
			return nil
		})
	}
	err = wg.Wait()

	return &dm.PropertyAggIndexResp{List: diDatas}, err
}
