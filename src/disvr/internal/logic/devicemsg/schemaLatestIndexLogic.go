package devicemsglogic

import (
	"context"
	"encoding/json"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/devices"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/disvr/internal/domain/deviceMsg"

	"github.com/i-Things/things/src/disvr/internal/svc"
	"github.com/i-Things/things/src/disvr/pb/di"

	"github.com/zeromicro/go-zero/core/logx"
)

type SchemaLatestIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSchemaLatestIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SchemaLatestIndexLogic {
	return &SchemaLatestIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取设备数据信息
func (l *SchemaLatestIndexLogic) SchemaLatestIndex(in *di.SchemaLatestIndexReq) (*di.SchemaIndexResp, error) {
	var (
		diDatas []*di.SchemaIndex
		total   int
	)
	temp, err := l.svcCtx.SchemaRepo.GetSchemaModel(l.ctx, in.ProductID)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	dd := l.svcCtx.SchemaMsgRepo
	switch in.Method {
	case devices.PropertyMethod, "":
		dataIDs := in.DataID
		if len(dataIDs) == 0 {
			dataIDs = temp.Property.GetIDs()
		}
		total = len(dataIDs)
		for _, v := range dataIDs {
			dds, err := dd.GetPropertyDataByID(l.ctx,
				deviceMsg.FilterOpt{
					Page:       def.PageInfo2{Size: 1},
					ProductID:  in.ProductID,
					DeviceName: []string{},
					DataID:     v,
					Order:      def.OrderDesc})
			if err != nil {
				l.Errorf("%s.GetPropertyDataByID err=%v", utils.FuncName(), err)
				return nil, errors.System.AddDetail(err)
			}
			var diData di.SchemaIndex
			if len(dds) == 0 {
				diData = di.SchemaIndex{
					Timestamp: 0,
					DataID:    v,
				}
			} else {
				devData := dds[0]
				diData = di.SchemaIndex{
					Timestamp: devData.TimeStamp.UnixMilli(),
					DataID:    devData.ID,
				}
				var payload []byte
				if param, ok := devData.Param.(string); ok {
					payload = []byte(param)
				} else {
					payload, _ = json.Marshal(devData.Param)
				}
				diData.GetValue = string(payload)

			}
			diDatas = append(diDatas, &diData)
			l.Infof("%s.get data=%+v", utils.FuncName(), diData)
		}
	default:
		return nil, errors.NotRealize.AddDetailf("multi method not implemt:%v", in.Method)
	}
	return &di.SchemaIndexResp{
		Total: int64(total),
		List:  diDatas,
	}, nil
}
