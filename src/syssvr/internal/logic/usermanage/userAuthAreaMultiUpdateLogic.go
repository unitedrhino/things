package usermanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAuthAreaMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UaaDB *relationDB.UserAuthAreaRepo
	UapDB *relationDB.UserAuthProjectRepo
}

func NewUserAuthAreaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthAreaMultiUpdateLogic {
	return &UserAuthAreaMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UaaDB:  relationDB.NewUserAuthAreaRepo(ctx),
		UapDB:  relationDB.NewUserAuthProjectRepo(ctx),
	}
}

func (l *UserAuthAreaMultiUpdateLogic) UserAuthAreaMultiUpdate(in *sys.UserAreaMultiUpdateReq) (*sys.Response, error) {
	if in.UserID == 0 {
		return nil, errors.Parameter.AddDetail(in.UserID).WithMsg("用户ID参数必填")
	}
	if in.ProjectID == 0 {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("项目id参数必填")
	}
	project, err := l.UapDB.FindOne(l.ctx, in.UserID, in.ProjectID)
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
	}
	po, err := checkUser(l.ctx, in.UserID)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("检查用户出错")
	} else if po == nil {
		return nil, errors.Parameter.AddDetail(err).WithMsg("检查用户不存在")
	}
	areas := ToAuthAreaDos(in.Areas)
	err = l.UaaDB.MultiUpdate(l.ctx, in.UserID, in.ProjectID, areas)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("用户数据权限保存失败")
	}
	if len(areas) == 0 && project != nil { //如果把项目下所有区域权限取消了,则项目权限默认也取消
		l.UapDB.Delete(l.ctx, in.UserID, in.ProjectID)
		InitCacheUserAuthProject(l.ctx, in.UserID)
	}
	if len(areas) != 0 && project == nil {
		l.UapDB.Insert(l.ctx, &relationDB.SysUserAuthProject{UserID: in.UserID, ProjectID: in.ProjectID})
		InitCacheUserAuthProject(l.ctx, in.UserID)
	}
	//更新 用户数据权限 缓存
	err = caches.SetUserAuthArea(l.ctx, in.UserID, in.ProjectID, areas)
	if err != nil {
		return nil, errors.Database.AddDetail(in.UserID).WithMsg("用户数据权限缓存失败")
	}

	return &sys.Response{}, nil
}
