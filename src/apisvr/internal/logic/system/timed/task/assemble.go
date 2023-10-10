package task

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/timedschedulersvr/client/timedscheduler"
)

func ToTaskInfoPb(in *types.TimedTaskInfo) *timedscheduler.TaskInfo {
	return &timedscheduler.TaskInfo{
		Id:             in.ID,
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

func ToTaskInfoTypes(in *timedscheduler.TaskInfo) *types.TimedTaskInfo {
	return &types.TimedTaskInfo{
		ID:             in.Id,
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
