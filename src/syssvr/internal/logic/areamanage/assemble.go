package areamanagelogic

import (
	"github.com/i-Things/things/shared/def"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/domain/area"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func transPoArrToPbTree(rootAreaID int64, poArr []*relationDB.SysAreaInfo) *sys.AreaInfo {
	pbList := make([]*sys.AreaInfo, 0, len(poArr))
	for _, po := range poArr {
		pbList = append(pbList, transPoToPb(po))
	}
	return buildPbTree(rootAreaID, pbList)
}

func transPoToPb(po *relationDB.SysAreaInfo) *sys.AreaInfo {
	parentAreaID := po.ParentAreaID
	if parentAreaID == 0 {
		parentAreaID = def.RootNode
	}
	return &sys.AreaInfo{
		CreatedTime:  po.CreatedTime.Unix(),
		AreaID:       int64(po.AreaID),
		ParentAreaID: parentAreaID,
		ProjectID:    int64(po.ProjectID),
		AreaName:     po.AreaName,
		Position:     logic.ToSysPoint(po.Position),
		Desc:         utils.ToRpcNullString(po.Desc),
	}
}

func buildPbTree(rootAreaID int64, pbList []*sys.AreaInfo) *sys.AreaInfo {
	var root *sys.AreaInfo
	// 将所有节点按照 parentID 分组
	nodeMap := make(map[int64][]*sys.AreaInfo)
	for _, pbOne := range pbList {
		if pbOne.AreaID == rootAreaID { // 找到根节点
			root = pbOne
		}
		nodeMap[pbOne.ParentAreaID] = append(nodeMap[pbOne.ParentAreaID], pbOne)
	}

	if root == nil { // 未找到根节点
		root = transPoToPb(&area.RootNode)
	}

	// 递归生成子树
	buildPbSubtree(root, nodeMap)

	return root
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
