package devicemsglogic

import (
	"context"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/logic"
	"gitee.com/unitedrhino/things/share/devices"
	"gitee.com/unitedrhino/things/share/domain/deviceMsg/msgThing"
	"gitee.com/unitedrhino/things/share/domain/schema"

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
	dc := devices.Core{
		ProductID:  in.ProductID,
		DeviceName: in.DeviceName,
	}
	_, err := logic.SchemaAccess(l.ctx, l.svcCtx, def.AuthRead, dc, nil)
	if err != nil {
		return nil, err
	}
	temp, err := l.svcCtx.DeviceSchemaRepo.GetData(l.ctx, dc)
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
	var lastBind int64
	uc := ctxs.GetUserCtxNoNil(l.ctx)
	if !uc.IsAdmin {
		di, err := l.svcCtx.DeviceCache.GetData(l.ctx, dc)
		if err != nil {
			return nil, err
		}
		lastBind = di.LastBind
	}
	for k, v := range dataMap {
		property := v
		dataID := k
		wait.Add(1)
		utils.Go(l.ctx, func() {
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
			if data != nil && lastBind != 0 {
				if data.TimeStamp.Before(time.Unix(lastBind, 0)) {
					data = nil
				}
			}
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
				sdef := property.Define
				if sdef.Type == schema.DataTypeArray {
					sdef = *sdef.ArrayInfo
				}
				if sdef.Type == schema.DataTypeStruct {
					dd, _ := schema.ParseDataID(dataID)
					if dd != nil && dd.Column != "" {
						sdef = sdef.Spec[dd.Column].DataType
					}
				}
				diData = dm.PropertyLogInfo{
					Timestamp: data.TimeStamp.UnixMilli(),
					DataID:    data.Identifier,
				}
				v, err := sdef.FmtValue(data.Param)
				if err == nil {
					diData.Value = utils.ToString(v)
				} else {
					diData.Value = utils.ToString(data.Param)
				}
			}
			diData.DataName = property.Name
			mutex.Lock()
			defer mutex.Unlock()
			diDatas = append(diDatas, &diData)
			l.Debugf("%s.get data=%+v", utils.FuncName(), diData)
		})
	}
	wait.Wait()
	return &dm.PropertyLogIndexResp{
		Total: int64(total),
		List:  diDatas,
	}, nil
}
