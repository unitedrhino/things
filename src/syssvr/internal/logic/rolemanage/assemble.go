package rolemanagelogic

import (
	"github.com/i-Things/things/src/syssvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToRoleInfoRpc(in *relationDB.SysRoleInfo) *sys.RoleInfo {
	if in == nil {
		return nil
	}
	return &sys.RoleInfo{
		Id:          in.ID,
		Name:        in.Name,
		Desc:        in.Desc,
		CreatedTime: in.CreatedTime.Unix(),
		Status:      in.Status,
	}
}
func ToRoleInfosRpc(in []*relationDB.SysRoleInfo) []*sys.RoleInfo {
	var ret []*sys.RoleInfo
	for _, v := range in {
		ret = append(ret, ToRoleInfoRpc(v))
	}
	return ret
}
