package area

import (
	"context"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
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

func (l *IndexLogic) Index(req *types.DataAreaIndexReq) (resp *types.DataAreaIndexResp, err error) {
	dto := &sys.DataAreaIndexReq{
		Page:      logic.ToSysPageRpc(req.Page),
		UserID:    req.UserID,
		ProjectID: req.ProjectID,
	}
	dmResp, err := l.svcCtx.DataM.DataAreaIndex(l.ctx, dto)
	if err != nil {
		l.Errorf("%s.rpc.DataAreaIndex req=%v err=%+v", utils.FuncName(), req, err)
		return nil, err
	}
	if len(dmResp.List) == 0 {
		return &types.DataAreaIndexResp{}, nil
	}
	var areaIDs []int64
	for _, v := range dmResp.List {
		areaIDs = append(areaIDs, v.AreaID)
	}
	areaInfos, err := l.svcCtx.AreaM.AreaInfoIndex(l.ctx, &sys.AreaInfoIndexReq{AreaIDs: areaIDs})
	if err != nil {
		return nil, err
	}
	var areaMap = map[int64]*sys.AreaInfo{}
	for _, v := range areaInfos.List {
		areaMap[v.AreaID] = v
	}
	list := ToDataAreaDetail(dmResp.List, areaMap)
	return &types.DataAreaIndexResp{
		Total: dmResp.Total,
		List:  list,
	}, nil
}
