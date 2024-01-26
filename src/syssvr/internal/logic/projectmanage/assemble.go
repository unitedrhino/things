package projectmanagelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/logic"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func transPoToPb(po *relationDB.SysProjectInfo) *sys.ProjectInfo {
	pb := &sys.ProjectInfo{
		CreatedTime: po.CreatedTime.Unix(),
		ProjectID:   int64(po.ProjectID),
		ProjectName: po.ProjectName,
		AdminUserID: po.AdminUserID,
		Desc:        utils.ToRpcNullString(po.Desc),
		Position:    logic.ToSysPoint(po.Position),
	}
	return pb
}
