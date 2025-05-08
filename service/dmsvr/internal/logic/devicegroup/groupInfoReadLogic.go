package devicegrouplogic

import (
	"context"
	"gitee.com/unitedrhino/share/def"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/domain/deviceGroup"
	"gitee.com/unitedrhino/things/service/dmsvr/internal/repo/relationDB"

	"gitee.com/unitedrhino/things/service/dmsvr/internal/svc"
	"gitee.com/unitedrhino/things/service/dmsvr/pb/dm"

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
func (l *GroupInfoReadLogic) GroupInfoRead(in *dm.GroupInfoReadReq) (*dm.GroupInfo, error) {
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
		if in.Purpose == "" {
			in.Purpose = deviceGroup.DictDefault
		}
		//ret, err := l.svcCtx.DictM.DictDetailRead(l.ctx, &sys.DictDetailReadReq{
		//	DictCode: deviceGroup.DictCode,
		//	Value:    in.Purpose,
		//})
		//if err == nil {
		//	po.Name = ret.Value
		//}
	case def.NotClassified:
		po = &relationDB.DmGroupInfo{
			ID:   def.NotClassified,
			Name: "自定义",
		}
	default:
		po, err = l.GiDB.FindOneByFilter(l.ctx, relationDB.GroupInfoFilter{Purpose: in.Purpose, ID: in.Id, WithProduct: true})
		if err != nil {
			return nil, err
		}
	}
	if !in.WithChildren {
		return ToGroupInfoPb(l.ctx, l.svcCtx, po), nil
	}
	children, err := l.GiDB.FindByFilter(l.ctx,
		relationDB.GroupInfoFilter{Purpose: in.Purpose, IDPath: po.IDPath, WithProduct: true}, nil)
	if err != nil {
		return nil, err
	}
	var ret = ToGroupInfoPb(l.ctx, l.svcCtx, po)
	if children != nil {
		var idMap = map[int64][]*dm.GroupInfo{}
		for _, v := range children {
			idMap[v.ParentID] = append(idMap[v.ParentID], ToGroupInfoPb(l.ctx, l.svcCtx, v))
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
