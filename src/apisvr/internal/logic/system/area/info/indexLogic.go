package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type IndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IndexLogic {
	return &IndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IndexLogic) Index(req *types.AreaInfoIndexReq) (resp *types.AreaInfoIndexResp, err error) {
	dmReq := &sys.AreaInfoIndexReq{
		Page:         logic.ToSysPageRpc(req.Page),
		ProjectID:    req.ProjectID,
		AreaIDs:      req.AreaIDs,
		ParentAreaID: req.ParentAreaID,
	}
	dmResp, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, dmReq)
	if err != nil {
		er := errors.Fmt(err)
		l.Errorf("%s.rpc.AreaManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}

	list := make([]*types.AreaInfo, 0, len(dmResp.List))
	for _, pb := range dmResp.List {
		var deviceCount *types.DeviceInfoCount
		if req.WithDeviceInfoCount {
			ret, err := l.svcCtx.DeviceM.DeviceInfoCount(l.ctx, &dm.DeviceInfoCountReq{
				TimeRange: nil,
				AreaIDs:   []int64{pb.AreaID},
				GroupIDs:  nil,
			})
			if err == nil {
				deviceCount = &types.DeviceInfoCount{
					Total:    ret.Total,
					Online:   ret.Online,
					Offline:  ret.Offline,
					Inactive: ret.Inactive,
					Unknown:  ret.Unknown,
				}
			}
		}

		list = append(list, ToAreaInfoTypes(pb, deviceCount))
	}

	return &types.AreaInfoIndexResp{
		Total: dmResp.Total,
		List:  list,
	}, nil
}
