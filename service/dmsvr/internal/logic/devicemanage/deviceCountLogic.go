package devicemanagelogic

import (
	"context"
	"fmt"
	"gitee.com/i-Things/core/service/syssvr/pb/sys"
	"gitee.com/i-Things/share/def"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProductInfoRepo
	DiDB *relationDB.DeviceInfoRepo
}

func NewDeviceCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceCountLogic {
	return &DeviceCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProductInfoRepo(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

const (
	RangeTypeAll = int64(iota)
	RangeTypeArea
	RangeTypeGroup
)

const (
	CountTypeStatus = "status"
	CountTypeType   = "type"
)

func (l *DeviceCountLogic) DeviceCount(in *dm.DeviceCountReq) (*dm.DeviceCountResp, error) {
	var list []*dm.DeviceCountInfo
	for _, v := range in.RangeIDs {
		f, err := l.FillFilter(in.RangeType, v)
		if err != nil {
			return nil, err
		}
		retMap, err := l.Count(in.CountTypes, *f)
		if err != nil {
			return nil, err
		}
		list = append(list, &dm.DeviceCountInfo{
			RangeID: v,
			Count:   retMap,
		})
	}
	return &dm.DeviceCountResp{List: list}, nil
}

func (l *DeviceCountLogic) FillFilter(rangeType int64, rangeID int64) (*relationDB.DeviceFilter, error) {
	var f relationDB.DeviceFilter
	switch rangeType {
	case RangeTypeAll:
		return &relationDB.DeviceFilter{}, nil
	case RangeTypeGroup:
		gds, err := relationDB.NewGroupDeviceRepo(l.ctx).FindByFilter(l.ctx, relationDB.GroupDeviceFilter{
			GroupIDs: []int64{rangeID},
		}, nil)
		if err != nil {
			return nil, err
		}
		if len(gds) == 0 {
			return nil, errors.NotFind
		}
		for _, v := range gds {
			f.Cores = append(f.Cores, &devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
		}
	case RangeTypeArea:
		switch rangeID {
		case def.RootNode:
			f.NotAreaIDs = []int64{def.NotClassified}
		case def.NotClassified:
			f.AreaIDs = []int64{def.NotClassified}
		default:
			ret, err := l.svcCtx.AreaM.AreaInfoRead(l.ctx, &sys.AreaInfoReadReq{AreaID: rangeID})
			if err != nil {
				return nil, err
			}
			f.AreaIDs = append(ret.ChildrenAreaIDs, ret.AreaID)
		}
	}
	return &f, nil
}

func (l *DeviceCountLogic) Count(countTypes []string, f relationDB.DeviceFilter) (map[string]int64, error) {
	var retMap = map[string]int64{}
	for _, countType := range countTypes {
		switch countType {
		case CountTypeType:
			// 获取 productID 统计
			productCount, err := l.DiDB.CountGroupByField(l.ctx, f, "product_id")
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
			retMap["typeDevice"] = deviceCount
			retMap["typeGateWay"] = gatewayCount
			retMap["typeSubSet"] = subsetCount
			retMap["typeUnknown"] = unknownCount
		case CountTypeStatus:
			diCount, err := l.DiDB.CountGroupByField(
				l.ctx, f, "is_online")
			if err != nil {
				if errors.Cmp(err, errors.NotFind) {
					return nil, errors.NotFind
				}
				return nil, err
			}
			onlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOnline)]
			offlineCount := diCount[fmt.Sprintf("%d", def.DeviceStatusOffline)]
			var allCount int64
			for _, v := range diCount {
				allCount += v
			}
			retMap["total"] = allCount
			retMap["online"] = onlineCount
			retMap["offline"] = offlineCount
		}
	}
	return retMap, nil
}
