package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/domain/deviceMsg/msgThing"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"sync"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type PropertyLogLatestIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyLogLatestIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogLatestIndexLogic {
	return &PropertyLogLatestIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *PropertyLogLatestIndexLogic) PropertyLogLatestIndex(in *dm.PropertyLogLatestIndexReq) (*dm.PropertyLogIndexResp, error) {
	var (
		diDatas []*dm.PropertyLogInfo
		total   int
	)
	temp, err := l.svcCtx.SchemaRepo.GetData(l.ctx, in.ProductID)
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
			var diData dm.PropertyLogInfo
			if data == nil {
				diData = dm.PropertyLogInfo{
					Timestamp: 0,
					DataID:    dataID,
				}
			} else {
				diData = dm.PropertyLogInfo{
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
	return &dm.PropertyLogIndexResp{
		Total: int64(total),
		List:  diDatas,
	}, nil
}
