package areamanagelogic

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func transPoArrToPbTree(root *relationDB.SysAreaInfo, poArr []*relationDB.SysAreaInfo) *sys.AreaInfo {
	pbList := make([]*sys.AreaInfo, 0, len(poArr))
	for _, po := range poArr {
		pbList = append(pbList, transPoToPb(po))
	}
	return buildPbTree(transPoToPb(root), pbList)
}

func buildPbTree(rootArea *sys.AreaInfo, pbList []*sys.AreaInfo) *sys.AreaInfo {
	// 将所有节点按照 parentID 分组
	nodeMap := make(map[int64][]*sys.AreaInfo)
	for _, pbOne := range pbList {
		nodeMap[pbOne.ParentAreaID] = append(nodeMap[pbOne.ParentAreaID], pbOne)
	}

	// 递归生成子树
	buildPbSubtree(rootArea, nodeMap)

	return rootArea
}

func transPoToPb(po *relationDB.SysAreaInfo) *sys.AreaInfo {
	parentAreaID := po.ParentAreaID
	if parentAreaID == 0 {
		parentAreaID = def.RootNode
	}
	return &sys.AreaInfo{
		CreatedTime:     po.CreatedTime.Unix(),
		AreaID:          int64(po.AreaID),
		ParentAreaID:    parentAreaID,
		ProjectID:       int64(po.ProjectID),
		AreaName:        po.AreaName,
		AreaNamePath:    GetNamePath(po.AreaNamePath),
		AreaIDPath:      GetIDPath(po.AreaIDPath),
		Position:        logic.ToSysPoint(po.Position),
		Desc:            utils.ToRpcNullString(po.Desc),
		LowerLevelCount: po.LowerLevelCount,
		ChildrenAreaIDs: po.ChildrenAreaIDs,
	}
}
func AreaInfosToPb(po []*relationDB.SysAreaInfo) (ret []*sys.AreaInfo) {
	for _, po := range po {
		ret = append(ret, transPoToPb(po))
	}
	return
}

func buildPbSubtree(node *sys.AreaInfo, nodeMap map[int64][]*sys.AreaInfo) {
	// 找到当前节点的子节点数组
	children := nodeMap[node.AreaID]

	// 递归生成子树
	for _, child := range children {
		buildPbSubtree(child, nodeMap)
	}

	// 将生成的子树数组作为当前节点的子节点数组
	node.Children = children
}
