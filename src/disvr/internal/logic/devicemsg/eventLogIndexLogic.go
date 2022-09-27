package devicemsglogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

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
func (l *EventLogIndexLogic) EventLogIndex(in *di.EventLogIndexReq) (*di.EventIndexResp, error) {
	var (
		diDatas []*di.EventIndex
		dd      = l.svcCtx.SchemaMsgRepo
		total   int64
	)
	dds, err := dd.GetEventDataByFilter(l.ctx, deviceMsg.FilterOpt{
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
		diData := di.EventIndex{
			Timestamp: devData.TimeStamp.UnixMilli(),
			Type:      devData.Type,
			DataID:    devData.ID,
		}
		var payload []byte
		payload, _ = json.Marshal(devData)
		diData.Params = string(payload)
		diDatas = append(diDatas, &diData)
		l.Infof("%s get data=%+v", utils.FuncName(), diData)
	}

	total, err = dd.GetEventCountByFilter(l.ctx, deviceMsg.FilterOpt{
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

	return &di.EventIndexResp{
		Total: total,
		List:  diDatas,
	}, nil
}
