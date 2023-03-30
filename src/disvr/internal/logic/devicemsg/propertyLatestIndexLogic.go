package devicemsglogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg/msgThing"
	"sync"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

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
func (l *PropertyLatestIndexLogic) PropertyLatestIndex(in *di.PropertyLatestIndexReq) (*di.PropertyIndexResp, error) {
	var (
		diDatas []*di.PropertyIndex
		total   int
	)
	temp, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	dd := l.svcCtx.SchemaMsgRepo
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
			var diData di.PropertyIndex
			if data == nil {
				diData = di.PropertyIndex{
					Timestamp: 0,
					DataID:    dataID,
				}
			} else {
				diData = di.PropertyIndex{
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
	return &di.PropertyIndexResp{
		Total: int64(total),
		List:  diDatas,
	}, nil
}
