package projectmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoDeleteLogic {
	return &ProjectInfoDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 删除项目
func (l *ProjectInfoDeleteLogic) ProjectInfoDelete(in *sys.ProjectInfoDeleteReq) (*sys.Response, error) {
	if in.ProjectID == 0 {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("项目ID参数必填")
	}

	po, err := checkProject(l.ctx, in.ProjectID)
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("检查项目出错")
	} else if po == nil {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
	}

	err = l.AiDB.DeleteByFilter(l.ctx, relationDB.AreaInfoFilter{ProjectID: in.ProjectID})
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("删除项目区域出错")
	}

	err = l.PiDB.Delete(l.ctx, in.ProjectID)
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("删除项目出错")
	}

	return &sys.Response{}, nil
}
