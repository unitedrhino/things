package devicemanagelogic

import (
	"context"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/devices"
	"gitee.com/i-Things/share/stores"
	"gitee.com/i-Things/share/utils"
	"gitee.com/i-Things/things/service/dmsvr/internal/logic"
	"gitee.com/i-Things/things/service/dmsvr/internal/repo/relationDB"
	"gitee.com/i-Things/things/service/dmsvr/internal/svc"
	"gitee.com/i-Things/things/service/dmsvr/pb/dm"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
)

type DeviceInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	DiDB *relationDB.DeviceInfoRepo
}

func NewDeviceInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeviceInfoIndexLogic {
	return &DeviceInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		DiDB:   relationDB.NewDeviceInfoRepo(ctx),
	}
}

// 获取设备信息列表
func (l *DeviceInfoIndexLogic) DeviceInfoIndex(in *dm.DeviceInfoIndexReq) (*dm.DeviceInfoIndexResp, error) {
	l.Infof("%s req=%v", utils.FuncName(), utils.Fmt(in))
	var (
		info  []*dm.DeviceInfo
		size  int64
		err   error
		cores []*devices.Core
	)
	l.ctx = ctxs.WithDefaultAllProject(l.ctx)
	if len(in.Devices) != 0 {
		for _, v := range in.Devices {
			cores = append(cores, &devices.Core{
				ProductID:  v.ProductID,
				DeviceName: v.DeviceName,
			})
		}
	}

	filter := relationDB.DeviceFilter{
		TenantCode:        in.TenantCode,
		ProductID:         in.ProductID,
		ProductIDs:        in.ProductIDs,
		AreaIDs:           in.AreaIDs,
		DeviceName:        in.DeviceName,
		DeviceNames:       in.DeviceNames,
		Gateway:           utils.Copy[devices.Core](in.Gateway),
		Cores:             cores,
		Iccid:             in.Iccid,
		Tags:              in.Tags,
		Range:             in.Range,
		Position:          logic.ToStorePoint(in.Position),
		DeviceAlias:       in.DeviceAlias,
		IsOnline:          in.IsOnline,
		ProductCategoryID: in.ProductCategoryID,
		Versions:          in.Versions,
		SharedType:        in.WithShared,
		CollectType:       in.WithCollect,
		DeviceType:        in.DeviceType,
		DeviceTypes:       in.DeviceTypes,
		GroupID:           in.GroupID,
		Status:            in.Status,
		NotGroupID:        in.NotGroupID,
		NotAreaID:         in.NotAreaID,
		UserID:            in.UserID,
		NetType:           in.NetType,
		HasOwner:          in.HasOwner,
		Distributor:       utils.Copy[stores.IDPathFilter](in.Distributor),
	}
	if in.RatedPower != nil {
		filter.RatedPower = stores.GetCmp(in.RatedPower.CmpType, in.RatedPower.Value)
	}
	if in.ExpTime != nil {
		filter.ExpTime = stores.GetCmp(in.ExpTime.CmpType, cast.ToTime(in.ExpTime.Value))
	}
	if err := ctxs.IsRoot(l.ctx); err == nil { //default租户才可以查看其他租户的设备
		ctxs.GetUserCtx(l.ctx).AllTenant = true
		defer func() {
			ctxs.GetUserCtx(l.ctx).AllTenant = false
		}()
	}
	size, err = l.DiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.DiDB.FindByFilter(l.ctx, filter,
		logic.ToPageInfo(in.Page).WithDefaultOrder(stores.OrderBy{Field: "is_online", Sort: stores.OrderAsc},
			stores.OrderBy{Field: "created_time", Sort: stores.OrderDesc}),
	)
	if err != nil {
		return nil, err
	}

	info = make([]*dm.DeviceInfo, 0, len(di))
	for _, v := range di {
		info = append(info, logic.ToDeviceInfo(l.ctx, l.svcCtx, v))
	}

	return &dm.DeviceInfoIndexResp{List: info, Total: size}, nil
}
