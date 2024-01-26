package info

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadLogic {
	return &ReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadLogic) Read(req *types.AreaInfoReadReq) (resp *types.AreaInfo, err error) {
	dmResp, err := l.svcCtx.AreaM.AreaInfoRead(l.ctx, &sys.AreaInfoReadReq{AreaID: req.AreaID, ProjectID: req.ProjectID, IsRetTree: req.IsRetTree})
	if er := errors.Fmt(err); er != nil {
		l.Errorf("%s.rpc.AreaManage req=%v err=%+v", utils.FuncName(), req, er)
		return nil, er
	}
	var deviceCount *types.DeviceInfoCount
	if req.WithDeviceInfoCount {
		ret, err := l.svcCtx.DeviceM.DeviceInfoCount(l.ctx, &dm.DeviceInfoCountReq{
			TimeRange: nil,
			AreaIDs:   append(dmResp.ChildrenAreaIDs, dmResp.AreaID),
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
	return ToAreaInfoTypes(dmResp, deviceCount), nil
}
