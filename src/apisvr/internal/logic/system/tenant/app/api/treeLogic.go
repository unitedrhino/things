package api

import (
	"context"
	"fmt"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TreeLogic {
	return &TreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func GenApiGroup(in []*sys.TenantApiInfo) []*types.TenantAppApiGroupInfo {
	var retMap = map[string][]*types.TenantApiInfo{}
	for _, v := range in {
		retMap[v.Info.Group] = append(retMap[v.Info.Group], ToTenantAppApiTypes(v))
	}
	var retList []*types.TenantAppApiGroupInfo
	var groupID int64
	for k, v := range retMap {
		groupID++
		retList = append(retList, &types.TenantAppApiGroupInfo{
			ID:       fmt.Sprintf("group%d", groupID),
			Name:     k,
			Children: v,
		})
	}
	return retList
}

func (l *TreeLogic) Tree(req *types.TenantAppApiIndexReq) (resp *types.TenantAppApiTreeResp, err error) {
	ret, err := l.svcCtx.TenantRpc.TenantAppApiIndex(l.ctx, &sys.TenantAppApiIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		Code:       req.Code,
		AppCode:    req.AppCode,
		ModuleCode: req.ModuleCode,
	})
	return &types.TenantAppApiTreeResp{
		List: GenApiGroup(ret.List),
	}, nil
}
