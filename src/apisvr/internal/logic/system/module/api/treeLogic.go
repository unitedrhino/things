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

func GenApiGroup(in []*sys.ApiInfo) []*types.ApiGroupInfo {
	var retMap = map[string][]*types.ApiInfo{}
	for _, v := range in {
		retMap[v.Group] = append(retMap[v.Group], ToApiInfoTypes(v))
	}
	var retList []*types.ApiGroupInfo
	var groupID int64
	for k, v := range retMap {
		groupID++
		retList = append(retList, &types.ApiGroupInfo{
			ID:       fmt.Sprintf("group%d", groupID),
			Name:     k,
			Children: v,
		})
	}
	return retList
}

func (l *TreeLogic) Tree(req *types.ApiInfoIndexReq) (resp *types.ApiInfoTreeResp, err error) {
	info, err := l.svcCtx.ModuleRpc.ModuleApiIndex(l.ctx, &sys.ApiInfoIndexReq{
		Page:       logic.ToSysPageRpc(req.Page),
		Route:      req.Route,
		Method:     req.Method,
		Group:      req.Group,
		Name:       req.Name,
		IsNeedAuth: req.IsNeedAuth,
	})
	if err != nil {
		return nil, err
	}

	return &types.ApiInfoTreeResp{List: GenApiGroup(info.List)}, nil
}
