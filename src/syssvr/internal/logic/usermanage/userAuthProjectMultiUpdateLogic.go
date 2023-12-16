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

type UserAuthProjectMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UapDB *relationDB.UserAuthProjectRepo
}

func NewUserAuthProjectMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthProjectMultiUpdateLogic {
	return &UserAuthProjectMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UapDB:  relationDB.NewUserAuthProjectRepo(ctx),
	}
}

func (l *UserAuthProjectMultiUpdateLogic) UserAuthProjectMultiUpdate(in *sys.UserAuthProjectMultiUpdateReq) (*sys.Response, error) {
	if in.UserID == 0 {
		return nil, errors.Parameter.AddDetail(in.UserID).WithMsg("用户ID参数必填")
	}
	po, err := checkUser(l.ctx, in.UserID)
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("检查用户出错")
	} else if po == nil {
		return nil, errors.Parameter.AddDetail(err).WithMsg("检查用户不存在")
	}
	projects := ToAuthProjectDos(in.Projects)
	err = l.UapDB.MultiUpdate(l.ctx, in.UserID, projects)
	if err != nil {
		return nil, errors.Database.AddDetail(err).WithMsg("用户数据权限保存失败")
	}

	//更新 用户数据权限 缓存
	err = caches.SetUserAuthProject(l.ctx, in.UserID, projects)
	if err != nil {
		return nil, errors.Database.AddDetail(in.UserID).WithMsg("用户数据权限缓存失败")
	}
	return &sys.Response{}, nil
}
