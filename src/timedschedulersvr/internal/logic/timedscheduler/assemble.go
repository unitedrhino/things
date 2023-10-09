package timedschedulerlogic

import (
	"github.com/i-Things/things/shared/domain/task"
	"github.com/i-Things/things/src/timedschedulersvr/client/timedscheduler"
	"github.com/i-Things/things/src/timedschedulersvr/internal/repo/relationDB"
)

func ToTimedTaskPo(in *timedscheduler.TaskInfo) *relationDB.TimedTask {
	if in == nil {
		return nil
	}
	return &relationDB.TimedTask{
		Group:          in.Group,
		Type:           in.Type,
		SubType:        in.SubType,
		Name:           in.Name,
		Code:           in.Code,
		Params:         in.Params,
		CronExpression: in.CronExpression,
		Status:         in.Status,
		EntryID:        in.EntryID,
		Priority:       in.Priority,
	}
}

func ToTimedTaskPb(in *relationDB.TimedTask) *timedscheduler.TaskInfo {
	if in == nil {
		return nil
	}
	return &timedscheduler.TaskInfo{
		Group:          in.Group,
		Type:           in.Type,
		SubType:        in.SubType,
		Name:           in.Name,
		Code:           in.Code,
		Params:         in.Params,
		CronExpression: in.CronExpression,
		Status:         in.Status,
		EntryID:        in.EntryID,
		Priority:       in.Priority,
	}
}

func ToTaskDo(in *timedscheduler.TaskInfo) *task.Info {
	if in == nil {
		return nil
	}
	return &task.Info{
		Group:    in.Group,
		Type:     in.Type,
		SubType:  in.SubType,
		Name:     in.Name,
		Code:     in.Code,
		Params:   in.Params,
		Priority: in.Priority,
	}
}

func PoToTaskDo(in *relationDB.TimedTask) *task.Info {
	if in == nil {
		return nil
	}
	return &task.Info{
		Group:    in.Group,
		Type:     in.Type,
		SubType:  in.SubType,
		Name:     in.Name,
		Code:     in.Code,
		Params:   in.Params,
		Priority: in.Priority,
	}
}
