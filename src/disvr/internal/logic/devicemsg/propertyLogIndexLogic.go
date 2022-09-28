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

type PropertyLogIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPropertyLogIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PropertyLogIndexLogic {
	return &PropertyLogIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *PropertyLogIndexLogic) PropertyLogIndex(in *di.PropertyLogIndexReq) (*di.PropertyIndexResp, error) {
	var (
		diDatas []*di.PropertyIndex
		dd      = l.svcCtx.SchemaMsgRepo
		total   int64
	)
	if in.Interval != 0 && in.ArgFunc == "" {
		return nil, errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	dds, err := dd.GetPropertyDataByID(l.ctx, deviceMsg.FilterOpt{
		Page: def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Page:      in.Page.GetPage(),
			Size:      in.Page.GetSize(),
		},
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		DataID:      in.DataID,
		Interval:    in.Interval,
		ArgFunc:     in.ArgFunc})
	if err != nil {
		l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
		return nil, err
	}
	for _, devData := range dds {
		diData := di.PropertyIndex{
			Timestamp: devData.TimeStamp.UnixMilli(),
			DataID:    devData.ID,
		}
		var payload []byte
		if param, ok := devData.Param.(string); ok {
			payload = []byte(param)
		} else {
			payload, _ = json.Marshal(devData.Param)
		}
		diData.Value = string(payload)
		diDatas = append(diDatas, &diData)
	}
	if in.ArgFunc == "" && in.Interval == 0 {
		total, err = dd.GetPropertyCountByID(l.ctx, deviceMsg.FilterOpt{
			Page: def.PageInfo2{
				TimeStart: in.TimeStart,
				TimeEnd:   in.TimeEnd,
				Page:      in.Page.GetPage(),
				Size:      in.Page.GetSize(),
			},
			ProductID:   in.ProductID,
			DeviceNames: in.DeviceNames,
			DataID:      in.DataID,
			Interval:    in.Interval,
			ArgFunc:     in.ArgFunc})
		if err != nil {
			l.Errorf("%s.GetPropertyCountByID err=%v", utils.FuncName(), err)
			return nil, err
		}
	}

	return &di.PropertyIndexResp{
		Total: total,
		List:  diDatas,
	}, nil
}
