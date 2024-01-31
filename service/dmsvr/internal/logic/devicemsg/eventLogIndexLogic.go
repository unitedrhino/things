package devicemsglogic

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/i-Things/things/service/dmsvr/internal/domain/deviceMsg/msgThing"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type EventLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEventLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventLogIndexLogic {
	return &EventLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *EventLogIndexLogic) EventLogIndex(in *dm.EventLogIndexReq) (*dm.EventIndexResp, error) {
	var (
		diDatas []*dm.EventIndex
		dd      = l.svcCtx.SchemaManaRepo
		total   int64
	)
	dds, err := dd.GetEventDataByFilter(l.ctx, msgThing.FilterOpt{
		Page: def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Page:      in.Page.GetPage(),
			Size:      in.Page.GetSize(),
		},
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		DataID:      in.DataID})
	if err != nil {
		l.Errorf("%s.GetEventDataByFilter err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}
	for _, devData := range dds {
		diData := dm.EventIndex{
			Timestamp: devData.TimeStamp.UnixMilli(),
			Type:      devData.Type,
			DataID:    devData.Identifier,
		}
		var payload []byte
		payload, _ = json.Marshal(devData.Params)
		diData.Params = string(payload)
		diDatas = append(diDatas, &diData)
		l.Infof("%s get data=%+v", utils.FuncName(), diData)
	}

	total, err = dd.GetEventCountByFilter(l.ctx, msgThing.FilterOpt{
		Page: def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Page:      in.Page.GetPage(),
			Size:      in.Page.GetSize(),
		},
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		DataID:      in.DataID})
	if err != nil {
		l.Errorf("%s.GetEventCountByFilter err=%v", utils.FuncName(), err)
		return nil, errors.System.AddDetail(err)
	}

	return &dm.EventIndexResp{
		Total: total,
		List:  diDatas,
	}, nil
}
