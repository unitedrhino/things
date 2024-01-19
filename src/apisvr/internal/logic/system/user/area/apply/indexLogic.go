package apply

import (
	"context"
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

func (l *IndexLogic) Index(req *types.UserAreaApplyIndexReq) (resp *types.UserAreaApplyIndexResp, err error) {
	ret, err := l.svcCtx.UserRpc.UserAreaApplyIndex(l.ctx, &sys.UserAreaApplyIndexReq{
		Page:      logic.ToSysPageRpc(req.Page),
		AuthTypes: req.AuthTypes,
	})
	if err != nil {
		return nil, err
	}
	var list []*types.UserAreaApplyInfo
	for _, v := range ret.List {
		list = append(list, &types.UserAreaApplyInfo{
			ID:          v.Id,
			AreaID:      v.AreaID,
			AuthType:    v.AuthType,
			CreatedTime: v.CreatedTime,
		})
	}
	return &types.UserAreaApplyIndexResp{
		Total: ret.Total,
		List:  list,
	}, nil
}
