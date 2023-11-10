package projectmanagelogic

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func transPoToPb(po *relationDB.SysProjectInfo) *sys.ProjectInfo {
	pb := &sys.ProjectInfo{
		CreatedTime: po.CreatedTime.Unix(),
		ProjectID:   int64(po.ProjectID),
		ProjectName: po.ProjectName,
		CompanyName: utils.ToRpcNullString(po.CompanyName),
		UserID:      po.UserID,
		Region:      utils.ToRpcNullString(po.Region),
		Address:     utils.ToRpcNullString(po.Address),
		Desc:        utils.ToRpcNullString(po.Desc),
	}
	return pb
}
