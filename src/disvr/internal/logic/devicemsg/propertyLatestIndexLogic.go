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
	for _, v := range dataIDs {
		dds, err := dd.GetPropertyDataByID(l.ctx,
			deviceMsg.FilterOpt{
				Page:        def.PageInfo2{Size: 1},
				ProductID:   in.ProductID,
				DeviceNames: []string{},
				DataID:      v,
				Order:       def.OrderDesc})
		if err != nil {
			l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
			return nil, errors.System.AddDetail(err)
		}
		var diData di.PropertyIndex
		if len(dds) == 0 {
			diData = di.PropertyIndex{
				Timestamp: 0,
				DataID:    v,
			}
		} else {
			devData := dds[0]
			diData = di.PropertyIndex{
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

		}
		diDatas = append(diDatas, &diData)
		l.Infof("%s.get data=%+v", utils.FuncName(), diData)
	}
	return &di.PropertyIndexResp{
		Total: int64(total),
		List:  diDatas,
	}, nil
}
