package devicegrouplogic

import (
	"context"
	"gitee.com/i-Things/share/def"
	"github.com/i-Things/things/service/dmsvr/internal/repo/relationDB"

	"github.com/i-Things/things/service/dmsvr/internal/svc"
	"github.com/i-Things/things/service/dmsvr/pb/dm"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
	GiDB *relationDB.GroupInfoRepo
}

func NewGroupInfoReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoReadLogic {
	return &GroupInfoReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		GiDB:   relationDB.NewGroupInfoRepo(ctx),
	}
}

// 获取分组信息详情
func (l *GroupInfoReadLogic) GroupInfoRead(in *dm.WithIDChildren) (*dm.GroupInfo, error) {
	var (
		po  *relationDB.DmGroupInfo
		err error
	)
	switch in.Id {
	case def.RootNode, 0:
		po = &relationDB.DmGroupInfo{
			ID:   def.RootNode,
			Name: "根节点",
		}
	case def.NotClassified:
		po = &relationDB.DmGroupInfo{
			ID:   def.NotClassified,
			Name: "自定义",
		}
	default:
		po, err = l.GiDB.FindOne(l.ctx, in.Id)
		if err != nil {
			return nil, err
		}
	}
	if !in.WithChildren {
		return ToGroupInfoPb(po), nil
	}
	children, err := l.GiDB.FindByFilter(l.ctx,
		relationDB.GroupInfoFilter{IDPath: po.IDPath}, nil)
	if err != nil {
		return nil, err
	}
	var ret = ToGroupInfoPb(po)
	if children != nil {
		var idMap = map[int64][]*dm.GroupInfo{}
		for _, v := range children {
			idMap[v.ParentID] = append(idMap[v.ParentID], ToGroupInfoPb(v))
		}
		fillDictInfoChildren(ret, idMap)
	}
	return ret, nil
}
func fillDictInfoChildren(node *dm.GroupInfo, nodeMap map[int64][]*dm.GroupInfo) {
	// 找到当前节点的子节点数组
	children := nodeMap[node.Id]
	for _, child := range children {
		fillDictInfoChildren(child, nodeMap)
	}
	node.Children = children
}
