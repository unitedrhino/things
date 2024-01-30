package datamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataProjectMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UapDB *relationDB.DataProjectRepo
}

func NewDataProjectMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataProjectMultiUpdateLogic {
	return &DataProjectMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UapDB:  relationDB.NewDataProjectRepo(ctx),
	}
}

func (l *DataProjectMultiUpdateLogic) DataProjectMultiUpdate(in *sys.DataProjectMultiUpdateReq) (*sys.Response, error) {
	if in.TargetID == 0 {
		return nil, errors.Parameter.AddDetail(in.TargetID).WithMsg("用户ID参数必填")
	}
	po, err := checkUser(l.ctx, in.TargetID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查用户出错")
	} else if po == nil {
		return nil, errors.Parameter.AddDetail(err).WithMsg("检查用户不存在")
	}
	projects := ToAuthProjectDos(in.Projects)
	err = l.UapDB.MultiUpdate(l.ctx, in.TargetID, projects)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("用户数据权限保存失败")
	}

	//更新 用户数据权限 缓存
	err = caches.SetUserAuthProject(l.ctx, in.TargetID, projects)
	if err != nil {
		return nil, errors.Database.AddDetail(in.TargetID).WithMsg("用户数据权限缓存失败")
	}
	return &sys.Response{}, nil
}
