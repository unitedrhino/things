package devicemanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/dmsvr/internal/svc"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
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
		info []*dm.DeviceInfo
		size int64
		err  error
	)

	filter := relationDB.DeviceFilter{
		ProductID:   in.ProductID,
		AreaIDs:     in.AreaIDs,
		DeviceName:  in.DeviceName,
		DeviceNames: in.DeviceNames,
		Tags:        in.Tags,
		Range:       in.Range,
		Position:    logic.ToStorePoint(in.Position),
		DeviceAlias: in.DeviceAlias,
		IsOnline:    in.IsOnline,
	}

	size, err = l.DiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}
	di, err := l.DiDB.FindByFilter(l.ctx, filter,
		logic.ToPageInfoWithDefault(in.Page, logic.ToPageInfo(in.Page,
			def.OrderBy{Filed: "created_time", Sort: def.OrderDesc},
			def.OrderBy{Filed: "product_id", Sort: def.OrderDesc})),
	)
	if err != nil {
		return nil, err
	}

	info = make([]*dm.DeviceInfo, 0, len(di))
	for _, v := range di {
		info = append(info, ToDeviceInfo(v))
	}

	return &dm.DeviceInfoIndexResp{List: info, Total: size}, nil
}
