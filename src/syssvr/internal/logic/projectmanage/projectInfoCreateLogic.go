package projectmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoCreateLogic {
	return &ProjectInfoCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 新增项目
func (l *ProjectInfoCreateLogic) ProjectInfoCreate(in *sys.ProjectInfo) (*sys.ProjectWithID, error) {
	if in.ProjectName == "" {
		return nil, errors.Parameter
	}

	po := &relationDB.SysProjectInfo{
		ProjectID:   stores.ProjectID(l.svcCtx.ProjectID.GetSnowflakeId()),
		ProjectName: in.ProjectName,
		//CompanyName: utils.ToEmptyString(in.CompanyName),
		AdminUserID: in.AdminUserID,
		//Region:      utils.ToEmptyString(in.Region),
		//Address:     utils.ToEmptyString(in.Address),
		Position: logic.ToStorePoint(in.Position),
		Desc:     utils.ToEmptyString(in.Desc),
	}

	err := l.PiDB.Insert(l.ctx, po)
	if err != nil {
		l.Errorf("%s.Insert err=%+v", utils.FuncName(), err)
		return nil, err
	}

	return &sys.ProjectWithID{ProjectID: int64(po.ProjectID)}, nil
}
