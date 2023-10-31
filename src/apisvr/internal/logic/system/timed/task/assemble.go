package task

import (
	"github.com/i-Things/things/src/apisvr/internal/types"
	"github.com/i-Things/things/src/timed/timedjobsvr/client/timedmanage"
)

func ToSendDelayReqPb(in *types.TimedTaskSendDelayReq) *timedmanage.TaskSendDelayReq {
	ret := timedmanage.TaskSendDelayReq{GroupCode: in.GroupCode, Code: in.Code}
	if in.Option != nil {
		ret.Option = &timedmanage.TaskDelayOption{
			Priority:  in.Option.Priority,
			ProcessIn: in.Option.ProcessIn,
			ProcessAt: in.Option.ProcessAt,
			Timeout:   in.Option.Timeout,
			Deadline:  in.Option.Deadline,
		}
	}
	if in.ParamSql != nil {
		ret.ParamSql = &timedmanage.TaskDelaySql{ExecContent: in.ParamSql.ExecContent}
	}
	if in.ParamQueue != nil {
		ret.ParamQueue = &timedmanage.TaskDelayQueue{Topic: in.ParamQueue.Topic, Payload: in.ParamQueue.Payload}
	}
	return &ret
}
