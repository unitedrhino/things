package timer

import (
	"encoding/json"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/src/timed/internal/domain"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
)

func FillTaskInfoDo(do *domain.TaskInfo, po *relationDB.TimedTask) error {
	if po.Group == nil {
		return errors.Parameter.AddMsg("任务没有找到任务组")
	}
	do.Code = po.Code
	do.GroupType = po.Group.Type
	do.GroupSubType = po.Group.SubType
	do.GroupCode = po.GroupCode
	if do.Params == "" { //如果没有传,则用数据库里的
		do.Params = po.Params
	}
	do.Env = po.Group.Env
	switch po.Group.Type {
	case domain.TaskGroupTypeQueue:
		var param domain.ParamQueue
		err := json.Unmarshal([]byte(do.Params), &param)
		if err != nil {
			return err
		}
		do.Queue = &param
	case domain.TaskGroupTypeSql:
		var sql domain.ParamSql
		err := json.Unmarshal([]byte(do.Params), &sql)
		if err != nil {
			return err
		}
		do.Sql = &domain.Sql{Param: sql}
		err = json.Unmarshal([]byte(po.Group.Config), &do.Sql.Config)
		if err != nil {
			return err
		}
	}
	return nil
}
