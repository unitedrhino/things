package role

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/syssvr/pb/sys"
)

func ToRoleInfoTypes(in *sys.RoleInfo) *types.RoleInfo {
	if in == nil {
		return nil
	}
	return &types.RoleInfo{
		ID:          in.Id,
		Name:        in.Name,
		Desc:        in.Desc,
		Status:      in.Status,
		CreatedTime: in.CreatedTime,
	}
}
func ToRoleInfosTypes(in []*sys.RoleInfo) (ret []*types.RoleInfo) {
	for _, v := range in {
		ret = append(ret, ToRoleInfoTypes(v))
	}
	return
}
