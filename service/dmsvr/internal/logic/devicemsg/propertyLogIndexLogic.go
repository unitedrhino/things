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
func (l *PropertyLogIndexLogic) PropertyLogIndex(in *dm.PropertyLogIndexReq) (*dm.PropertyIndexResp, error) {
	var (
		diDatas []*dm.PropertyIndex
		dd      = l.svcCtx.SchemaManaRepo
		total   int64
	)
	if in.Interval != 0 && in.ArgFunc == "" {
		return nil, errors.Parameter.AddMsg("填写了间隔就必须填写聚合函数")
	}
	dds, err := dd.GetPropertyDataByID(l.ctx, msgThing.FilterOpt{
		Page: def.PageInfo2{
			TimeStart: in.TimeStart,
			TimeEnd:   in.TimeEnd,
			Page:      in.Page.GetPage(),
			Size:      in.Page.GetSize(),
		},
		ProductID:   in.ProductID,
		DeviceNames: in.DeviceNames,
		Order:       in.Order,
		DataID:      in.DataID,
		Fill:        in.Fill,
		Interval:    in.Interval,
		ArgFunc:     in.ArgFunc})
	if err != nil {
		l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
		return nil, err
	}
	for _, devData := range dds {
		diData := dm.PropertyIndex{
			Timestamp: devData.TimeStamp.UnixMilli(),
			DataID:    devData.Identifier,
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
		total, err = dd.GetPropertyCountByID(l.ctx, msgThing.FilterOpt{
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

	return &dm.PropertyIndexResp{
		Total: total,
		List:  diDatas,
	}, nil
}
