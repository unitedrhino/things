package task

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
)

func ToSendDelayReqPb(in *types.TimedTaskSendReq) *timedmanage.TaskSendReq {
	ret := timedmanage.TaskSendReq{GroupCode: in.GroupCode, Code: in.Code}
	if in.Option != nil {
		ret.Option = &timedmanage.TaskSendOption{
			Priority:  in.Option.Priority,
			ProcessIn: in.Option.ProcessIn,
			ProcessAt: in.Option.ProcessAt,
			Timeout:   in.Option.Timeout,
			Deadline:  in.Option.Deadline,
		}
	}
	if in.ParamSql != nil {
		ret.ParamSql = &timedmanage.TaskParamSql{ExecContent: in.ParamSql.ExecContent, Param: in.ParamSql.Param}
	}
	if in.ParamQueue != nil {
		ret.ParamQueue = &timedmanage.TaskParamQueue{Topic: in.ParamQueue.Topic, Payload: in.ParamQueue.Payload}
	}
	return &ret
}

func ToGroupPb(in *types.TimedTaskGroup) *timedmanage.TaskGroup {
	if in == nil {
		return nil
	}
	return &timedmanage.TaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}

func ToGroupTypes(in *timedmanage.TaskGroup) *types.TimedTaskGroup {
	if in == nil {
		return nil
	}
	return &types.TimedTaskGroup{
		Code:     in.Code,
		Name:     in.Name,
		Type:     in.Type,
		SubType:  in.SubType,
		Priority: in.Priority,
		Env:      in.Env,
		Config:   in.Config,
	}
}
func ToTaskGroupsTypes(in []*timedmanage.TaskGroup) (ret []*types.TimedTaskGroup) {
	for _, v := range in {
		ret = append(ret, ToGroupTypes(v))
	}
	return
}

func ToTaskInfoPb(in *types.TimedTaskInfo) *timedmanage.TaskInfo {
	if in == nil {
		return nil
	}
	return &timedmanage.TaskInfo{
		GroupCode: in.GroupCode,
		Type:      in.Type,
		Name:      in.Name,
		Code:      in.Code,
		Params:    in.Params,
		CronExpr:  in.CronExpr,
		Status:    in.Status,
		Priority:  in.Priority,
	}
}

func ToTaskInfoTypes(in *timedmanage.TaskInfo) *types.TimedTaskInfo {
	if in == nil {
		return nil
	}
	return &types.TimedTaskInfo{
		GroupCode: in.GroupCode,
		Type:      in.Type,
		Name:      in.Name,
		Code:      in.Code,
		Params:    in.Params,
		CronExpr:  in.CronExpr,
		Status:    in.Status,
		Priority:  in.Priority,
	}
}

func ToTaskInfosTypes(in []*timedmanage.TaskInfo) (ret []*types.TimedTaskInfo) {
	for _, v := range in {
		ret = append(ret, ToTaskInfoTypes(v))
	}
	return
}
