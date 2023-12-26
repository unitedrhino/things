package modulemanagelogic

import (
	"context"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"

	"github.com/i-Things/things/src/syssvr/internal/svc"
	"github.com/i-Things/things/src/syssvr/pb/sys"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModuleMenuIndexLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewModuleMenuIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModuleMenuIndexLogic {
	return &ModuleMenuIndexLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ModuleMenuIndexLogic) ModuleMenuIndex(in *sys.MenuInfoIndexReq) (*sys.MenuInfoIndexResp, error) {
	f := relationDB.MenuInfoFilter{
		ModuleCode: in.ModuleCode,
		MenuIDs:    in.MenuIDs,
	}
	resp, err := relationDB.NewMenuInfoRepo(l.ctx).FindByFilter(l.ctx, f, nil)
	if err != nil {
		return nil, err
	}
	//total, err := relationDB.NewMenuInfoRepo(l.ctx).CountByFilter(l.ctx, f)
	//if err != nil {
	//	return nil, err
	//}
	info := make([]*sys.MenuInfo, 0, len(resp))
	if !in.IsRetTree {
		for _, v := range resp {
			i := logic.ToMenuInfoPb(v)
			info = append(info, i)
		}
		return &sys.MenuInfoIndexResp{List: info}, nil
	}

	var (
		pidMap = make(map[int64][]*sys.MenuInfo, len(resp))
		idMap  = make(map[int64]*sys.MenuInfo, len(resp))
	)
	for _, v := range resp {
		i := logic.ToMenuInfoPb(v)
		idMap[i.Id] = i
		if i.ParentID == 1 { //根节点
			info = append(info, i)
			continue
		}
		pidMap[i.ParentID] = append(pidMap[i.ParentID], i)
	}
	fillChildren(info, pidMap)
	return &sys.MenuInfoIndexResp{List: info}, nil
}
func fillChildren(in []*sys.MenuInfo, pidMap map[int64][]*sys.MenuInfo) {
	for _, v := range in {
		cs := pidMap[v.Id]
		if cs != nil {
			v.Children = cs
			fillChildren(cs, pidMap)
		}
	}
}
