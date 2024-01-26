package info

import (
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/apisvr/internal/logic"
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToProjectPb(in *types.ProjectInfo) *sys.ProjectInfo {
	if in == nil {
		return nil
	}
	return &sys.ProjectInfo{
		ProjectID:   in.ProjectID,
		ProjectName: in.ProjectName,
		AdminUserID: in.AdminUserID,
		Position:    logic.ToSysPointRpc(in.Position),
		Desc:        utils.ToRpcNullString(in.Desc),
	}
}
