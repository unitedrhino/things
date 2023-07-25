package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceTypeCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	DiDB *relationDB.DeviceInfoRepo
}

func NewDeviceTypeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceTypeCountLogic {
	return &DeviceTypeCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

// 设备类型
func (l *DeviceTypeCountLogic) DeviceTypeCount(in *dm.DeviceTypeCountReq) (*dm.DeviceTypeCountResp, error) {
	// 获取 productID 统计
	productCount, err := l.DiDB.CountGroupByField(
		l.ctx,
		relationDB.DeviceFilter{
			LastLoginTime: struct {
				Start int64
				End   int64
			}{Start: in.StartTime, End: in.EndTime},
		},
		"product_id",
	)

	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind
		}
		return nil, err
	}
	productIDs := make([]string, 0, len(productCount))
	for productID := range productCount {
		productIDs = append(productIDs, productID)
	}

	// 通过 productID 查找 DeviceType
	productIDList, err := l.PiDB.FindByFilter(l.ctx, relationDB.ProductFilter{
		ProductIDs: productIDs,
	}, nil)

	if err != nil {
		if errors.Cmp(err, errors.NotFind) {
			return nil, errors.NotFind
		}
		return nil, err
	}
	// 计算
	productMap := make(map[string]int64, 0)
	for _, v := range productIDList {
		productMap[v.ProductID] = v.DeviceType
	}

	var deviceCount, gatewayCount, subsetCount, unknownCount int64
	for productID, v := range productCount {
		productType := productMap[productID]
		switch productType {
		case def.DeviceTypeDevice:
			deviceCount += v
		case def.DeviceTypeGateway:
			gatewayCount += v
		case def.DeviceTypeSubset:
			subsetCount += v
		default:
			unknownCount += v
		}
	}

	return &dm.DeviceTypeCountResp{
		Device:  deviceCount,
		Gateway: gatewayCount,
		Subset:  subsetCount,
		Unknown: unknownCount,
	}, nil
}
