package projectmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoIndexLogic {
	return &ProjectInfoIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 获取项目信息列表
func (l *ProjectInfoIndexLogic) ProjectInfoIndex(in *sys.ProjectInfoIndexReq) (*sys.ProjectInfoIndexResp, error) {
	var (
		list  []*sys.ProjectInfo
		total int64
		err   error
	)

	filter := relationDB.ProjectInfoFilter{
		ProjectIDs:  in.ProjectIDs,
		ProjectName: in.ProjectName,
	}
	total, err = l.PiDB.CountByFilter(l.ctx, filter)
	if err != nil {
		return nil, err
	}

	poArr, err := l.PiDB.FindByFilter(l.ctx, filter, logic.ToPageInfo(in.Page))
	if err != nil {
		return nil, err
	}

	list = make([]*sys.ProjectInfo, 0, len(poArr))
	for _, po := range poArr {
		list = append(list, transPoToPb(po))
	}
	return &sys.ProjectInfoIndexResp{List: list, Total: total}, nil
}
