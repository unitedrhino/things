package areamanagelogic

import (
	"context"
	"github.com/i-Things/things/shared/ctxs"
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
	"github.com/samber/lo"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AreaInfoTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	AiDB *relationDB.AreaInfoRepo
}

func NewAreaInfoTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreaInfoTreeLogic {
	return &AreaInfoTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		AiDB:   relationDB.NewAreaInfoRepo(ctx),
	}
}

// 获取区域信息树
func (l *AreaInfoTreeLogic) AreaInfoTree(in *sys.AreaInfoTreeReq) (*sys.AreaInfoTreeResp, error) {
	var (
		err   error
		poArr []*relationDB.SysAreaInfo
		root  = int64(def.RootNode)
	)

	if in.AreaID != 0 && in.AreaID != def.RootNode { //如果传的是root节点,也会返回所有区域,如果传的是未分类的id,则也会把未分类下的所有设备加入
		root = in.AreaID
		areaIDs, err := l.AiDB.FindIDsWithChildren(l.ctx, in.AreaID)
		if err != nil {
			return nil, err
		}
		poArr, err = l.AiDB.FindByFilter(l.ctx,
			relationDB.AreaInfoFilter{AreaIDs: areaIDs}, nil)
		if err != nil {
			return nil, err
		}
	} else {
		poArr, err = l.AiDB.FindByFilter(l.ctx,
			relationDB.AreaInfoFilter{ProjectID: in.ProjectID}, nil)
		if err != nil {
			return nil, err
		}
	}
	ctxs.SetInnerCtx(l.ctx, ctxs.InnerCtx{AllData: true})
	poArr = l.checkMissingParentIdMenuIndex(poArr)
	return &sys.AreaInfoTreeResp{Tree: transPoArrToPbTree(root, poArr)}, nil
}
func (l *AreaInfoTreeLogic) checkMissingParentIdMenuIndex(areaInfos []*relationDB.SysAreaInfo) []*relationDB.SysAreaInfo {
	missingParentIds := findMissingParentIds(areaInfos)
	if len(missingParentIds) > 0 {
		areaIDs := lo.Keys(missingParentIds)
		areaInfo, err := l.AiDB.FindByFilter(l.ctx, relationDB.AreaInfoFilter{AreaIDs: areaIDs}, nil)
		if err != nil {
			l.Errorf("MenuIndex find menu_info err,menuIds:%d,err:%v", areaIDs, err)
			return areaInfos
		}
		areaInfos = append(areaInfos, areaInfo...)
		return l.checkMissingParentIdMenuIndex(areaInfos) //多级嵌套需要遍历
	}
	return areaInfos
}

func findMissingParentIds(areaInfos []*relationDB.SysAreaInfo) map[int64]bool {
	missingParentIds := make(map[int64]bool)
	ids := make(map[int64]bool)
	for _, menu := range areaInfos {
		ids[int64(menu.AreaID)] = true
	}
	for _, menu := range areaInfos {
		if !ids[menu.ParentAreaID] && menu.ParentAreaID != 1 {
			missingParentIds[menu.ParentAreaID] = true
		}
	}
	return missingParentIds
}
