package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/devices"
	"gitee.com/unitedrhino/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/share/domain/schema"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"sync"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
		dataMap map[string]*schema.Property
	)
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}, nil)
	if err != nil {
		return nil, err
	}
	temp, err := l.svcCtx.SchemaRepo.GetData(l.ctx, in.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	dd := l.svcCtx.SchemaManaRepo
	if len(in.DataIDs) == 0 {
		dataMap = temp.Property.GetMap()
	} else {
		dataMap = temp.Property.GetMapWithIDs(in.DataIDs...)
	}
	total = len(dataMap)
	wait := sync.WaitGroup{}
	mutex := sync.Mutex{}
	for k, v := range dataMap {
		property := v
		dataID := k
		wait.Add(1)
		utils.Go(l.ctx, func() {
			func() {
				defer wait.Done()
				data, err := dd.GetLatestPropertyDataByID(l.ctx, property, msgThing.LatestFilter{
					ProductID:  in.ProductID,
					DeviceName: in.DeviceName,
					DataID:     dataID,
				})
				if err != nil {
					l.Errorf("%s.GetLatestPropertyDataByID err=%v", utils.FuncName(), utils.Fmt(err))
					return
				}
				var diData dm.PropertyLogInfo
				if data == nil {
					v, err := property.Define.GetDefaultValue()
					if err != nil {
						l.Errorf("%s.GetDefaultValue err=%v", utils.FuncName(), utils.Fmt(err))
						return
					}

					diData = dm.PropertyLogInfo{
						Timestamp: 0,
						DataID:    dataID,
						Value:     utils.ToString(v),
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
				l.Debugf("%s.get data=%+v", utils.FuncName(), diData)
			}()
		})
	}
	wait.Wait()
	return &dm.PropertyLogIndexResp{
		Total: int64(total),
		List:  diDatas,
	}, nil
}
