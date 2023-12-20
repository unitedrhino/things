package projectmanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProjectInfoUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	PiDB *relationDB.ProjectInfoRepo
}

func NewProjectInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProjectInfoUpdateLogic {
	return &ProjectInfoUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		PiDB:   relationDB.NewProjectInfoRepo(ctx),
	}
}

// 更新项目
func (l *ProjectInfoUpdateLogic) ProjectInfoUpdate(in *sys.ProjectInfo) (*sys.Response, error) {
	if in.ProjectID == 0 {
		return nil, errors.Parameter
	}

	po, err := checkProject(l.ctx, in.ProjectID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查项目出错")
	} else if po == nil {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("检查项目不存在")
	}

	l.setPoByPb(po, in)

	err = l.PiDB.Update(l.ctx, po)
	if err != nil {
		return nil, err
	}
	return &sys.Response{}, nil
}
func (l *ProjectInfoUpdateLogic) setPoByPb(po *relationDB.SysProjectInfo, pb *sys.ProjectInfo) {
	if pb.ProjectName != "" {
		po.ProjectName = pb.ProjectName
	}
	if pb.CompanyName != nil {
		po.CompanyName = pb.CompanyName.GetValue()
	}
	if pb.UserID != 0 {
		po.UserID = pb.UserID
	}
	if pb.Region != nil {
		po.Region = pb.Region.GetValue()
	}
	if pb.Address != nil {
		po.Address = pb.Address.GetValue()
	}
	if pb.Desc != nil {
		po.Desc = pb.Desc.GetValue()
	}
}
