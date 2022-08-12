package logic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/domain/service/deviceData"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataSchemaLatestIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDataSchemaLatestIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataSchemaLatestIndexLogic {
	return &DataSchemaLatestIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *DataSchemaLatestIndexLogic) DataSchemaLatestIndex(in *dm.DataSchemaLatestIndexReq) (*dm.DataSchemaIndexResp, error) {
	var (
		dmDatas []*dm.DataSchemaIndex
		total   int
	)
	temp, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	dd := l.svcCtx.DeviceDataRepo
	switch in.Method {
	case def.PropertyMethod, "":
		dataIDs := in.DataID
		if len(dataIDs) == 0 {
			dataIDs = temp.Properties.GetIDs()
		}
		total = len(dataIDs)
		for _, v := range dataIDs {
			dds, err := dd.GetPropertyDataByID(l.ctx,
				deviceData.FilterOpt{
					Page:       def.PageInfo2{Size: 1},
					ProductID:  in.ProductID,
					DeviceName: in.DeviceName,
					DataID:     v,
					Order:      "desc"})
			if err != nil {
				l.Errorf("HandleData|GetPropertyDataByID|err=%v", err)
				return nil, errors.System.AddDetail(err)
			}
			var dmData dm.DataSchemaIndex
			if len(dds) == 0 {
				dmData = dm.DataSchemaIndex{
					Timestamp: 0,
					DataID:    v,
				}
			} else {
				devData := dds[0]
				dmData = dm.DataSchemaIndex{
					Timestamp: devData.TimeStamp.UnixMilli(),
					DataID:    devData.ID,
				}
				var payload []byte
				if param, ok := devData.Param.(string); ok {
					payload = []byte(param)
				} else {
					payload, _ = json.Marshal(devData.Param)
				}
				dmData.GetValue = string(payload)

			}
			dmDatas = append(dmDatas, &dmData)
			l.Infof("GetDeviceLogLogic|get data=%+v", dmData)
		}
	default:
		return nil, errors.NotRealize.AddDetailf("multi method not implemt:%v", in.Method)
	}
	return &dm.DataSchemaIndexResp{
		Total: int64(total),
		List:  dmDatas,
	}, nil
}
