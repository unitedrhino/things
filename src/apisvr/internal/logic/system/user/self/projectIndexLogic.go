package self

import (
	"context"
	"github.com/i-Things/things/src/apisvr/internal/logic/system"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/apisvr/internal/svc"
	"github.com/i-Things/things/src/apisvr/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProjectIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectIndexLogic {
	return &ProjectIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProjectIndexLogic) ProjectIndex() (resp *types.ProjectInfoIndexResp, err error) {
	var (
		projectInfos []*sys.ProjectInfo
	)
	//uc := ctxs.GetUserCtx(l.ctx)
	//ret, err := l.svcCtx.UserRpc.UserProjectIndex(l.ctx, &sys.UserProjectIndexReq{
	//	UserID: uc.UserID,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//if len(ret.List) != 0 {
	//	var projectIDs []int64
	//	for _, v := range ret.List {
	//		projectIDs = append(projectIDs, v.ProjectID)
	//	}
	//	ret2, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, &sys.ProjectInfoIndexReq{ProjectIDs: projectIDs})
	//	if err != nil {
	//		return nil, err
	//	}
	//	projectInfos = ret2.List
	//}
	ret2, err := l.svcCtx.ProjectM.ProjectInfoIndex(l.ctx, &sys.ProjectInfoIndexReq{})
	if err != nil {
		return nil, err
	}
	projectInfos = ret2.List
	return &types.ProjectInfoIndexResp{
		List: system.ProjectInfosToApi(projectInfos),
	}, nil
}
