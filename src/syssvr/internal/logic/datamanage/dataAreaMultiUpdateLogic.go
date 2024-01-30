package datamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/caches"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type DataAreaMultiUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	UaaDB *relationDB.DataAreaRepo
	UapDB *relationDB.DataProjectRepo
}

func NewDataAreaMultiUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DataAreaMultiUpdateLogic {
	return &DataAreaMultiUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		UaaDB:  relationDB.NewDataAreaRepo(ctx),
		UapDB:  relationDB.NewDataProjectRepo(ctx),
	}
}

func (l *DataAreaMultiUpdateLogic) DataAreaMultiUpdate(in *sys.DataAreaMultiUpdateReq) (*sys.Response, error) {
	if in.TargetID == 0 {
		return nil, errors.Parameter.AddDetail(in.TargetID).WithMsg("TargetID参数必填")
	}
	if in.ProjectID == 0 {
		return nil, errors.Parameter.AddDetail(in.ProjectID).WithMsg("项目id参数必填")
	}
	project, err := l.UapDB.FindOne(l.ctx, in.TargetType, in.TargetID, in.ProjectID)
	if err != nil {
		if !errors.Cmp(err, errors.NotFind) {
			return nil, err
		}
	}
	//po, err := checkUser(l.ctx, in.TargetID)
	//if err != nil {
	//	return nil, errors.Fmt(err).WithMsg("检查用户出错")
	//} else if po == nil {
	//	return nil, errors.Parameter.AddDetail(err).WithMsg("检查用户不存在")
	//}
	areas := ToAuthAreaDos(in.Areas)
	err = l.UaaDB.MultiUpdate(l.ctx, in.TargetID, in.ProjectID, areas)
	if err != nil {
		return nil, errors.Fmt(err).WithMsg("用户数据权限保存失败")
	}
	if len(areas) == 0 && project != nil { //如果把项目下所有区域权限取消了,则项目权限默认也取消
		l.UapDB.Delete(l.ctx, in.TargetType, in.TargetID, in.ProjectID)
		InitCacheUserAuthProject(l.ctx, in.TargetID)
	}
	if len(areas) != 0 && project == nil {
		l.UapDB.Insert(l.ctx, &relationDB.SysDataProject{TargetType: def.TargetUser, TargetID: in.TargetID, ProjectID: in.ProjectID})
		InitCacheUserAuthProject(l.ctx, in.TargetID)
	}
	//更新 用户数据权限 缓存
	err = caches.SetUserAuthArea(l.ctx, in.TargetID, in.ProjectID, areas)
	if err != nil {
		return nil, errors.Database.AddDetail(in.TargetID).WithMsg("用户数据权限缓存失败")
	}

	return &sys.Response{}, nil
}
