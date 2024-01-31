package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/core/shared/errors"
	"gitee.com/i-Things/core/shared/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/service/dmsvr/pb/dm"
	"sync"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLatestIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyLatestIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLatestIndexLogic {
	return &PropertyLatestIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *PropertyLatestIndexLogic) PropertyLatestIndex(in *dm.PropertyLatestIndexReq) (*dm.PropertyIndexResp, error) {
	var (
		diDatas []*dm.PropertyIndex
		total   int
	)
	temp, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	dd := l.svcCtx.SchemaManaRepo
	dataIDs := in.DataIDs
	if len(dataIDs) == 0 {
		dataIDs = temp.Property.GetIDs()
	}
	total = len(dataIDs)
	wait := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for _, v := range dataIDs {
		wait.Add(1)
		go func(dataID string) {
			defer wait.Done()
			data, err := dd.GetLatestPropertyDataByID(l.ctx, msgThing.LatestFilter{
				ProductID:  in.ProductID,
				DeviceName: in.DeviceName,
				DataID:     dataID,
			})
			if err != nil {
				l.Errorf("%s.GetLatestPropertyDataByID err=%v", utils.FuncName(), err)
				return
			}
			var diData dm.PropertyIndex
			if data == nil {
				diData = dm.PropertyIndex{
					Timestamp: 0,
					DataID:    dataID,
				}
			} else {
				diData = dm.PropertyIndex{
					Timestamp: data.TimeStamp.UnixMilli(),
					DataID:    data.Identifier,
				}
				var payload []byte
				if param, ok := data.Param.(string); ok {
					payload = []byte(param)
				} else {
					payload, _ = json.Marshal(data.Param)
				}
				diData.Value = string(payload)

			}
			mutex.Lock()
			defer mutex.Unlock()
			diDatas = append(diDatas, &diData)
			l.Infof("%s.get data=%+v", utils.FuncName(), diData)
		}(v)
	}
	wait.Wait()
	return &dm.PropertyIndexResp{
		Total: int64(total),
		List:  diDatas,
	}, nil
}
