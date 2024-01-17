package info

import (
	"context"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/dmsvr/pb/dm"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.AreaInfo) (*types.AreaWithID, error) {
	if req.ParentAreaID == 0 {
		req.ParentAreaID = def.RootNode
	}
	if req.ParentAreaID != def.RootNode {
		dmRep, err := l.svcCtx.DeviceM.DeviceInfoIndex(l.ctx, &dm.DeviceInfoIndexReq{
			Page:    &dm.PageInfo{Page: 1, Size: 2}, //只需要知道是否有设备即可
			AreaIDs: []int64{req.ParentAreaID}})
		if err != nil {
			return nil, err
		}
		if len(dmRep.List) != 0 {
			return nil, errors.Parameter.AddMsg("父级区域已绑定了设备，不允许再添加子区域")
		}
	}

	dmReq := &sys.AreaInfo{
		ParentAreaID: req.ParentAreaID,
		ProjectID:    req.ProjectID,
		AreaName:     req.AreaName,
		Position:     logic.ToSysPointRpc(req.Position),
		Desc:         utils.ToRpcNullString(req.Desc),
	}
	resp, err := l.svcCtx.AreaM.AreaInfoCreate(l.ctx, dmReq)
	if er := errors.Fmt(err); er != nil {
		l.Errorf("%s.rpc.AreaManage req=%v err=%v", utils.FuncName(), req, er)
		return nil, er
	}
	return &types.AreaWithID{AreaID: resp.AreaID}, nil
}
