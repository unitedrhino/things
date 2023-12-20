package projectmanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoReadLogic {
	return &ProjectInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 获取项目信息详情
func (l *ProjectInfoReadLogic) ProjectInfoRead(in *sys.ProjectWithID) (*sys.ProjectInfo, error) {
	po, err := l.PiDB.FindOne(l.ctx, in.ProjectID)
	if err != nil {
		return nil, err
	}
	return transPoToPb(po), nil
}
