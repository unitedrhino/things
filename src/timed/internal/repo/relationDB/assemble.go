package relationDB

import (
	"encoding/json"
	"github.com/i-Things/things/src/timed/internal/domain"
)

func ToTaskInfoDo(po *TimedTask) *domain.TaskInfo {
	var do domain.TaskInfo
	do.Code = po.Code
	do.GroupType = po.Group.Type
	do.GroupSubType = po.Group.SubType
	do.GroupCode = po.GroupCode
	do.Env = po.Group.Env
	switch po.Group.Type {
	case domain.TaskGroupTypeQueue:
		var param domain.ParamQueue
		json.Unmarshal([]byte(po.Params), &param)
		do.Queue = &param
	case domain.TaskGroupTypeSql:
		var sql domain.ParamSql
		json.Unmarshal([]byte(po.Params), &sql)
		do.Sql = &domain.Sql{Param: sql}
		json.Unmarshal([]byte(po.Group.Config), &do.Sql.Config)
	}
	return &do
}
