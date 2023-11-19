package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/i-Things/things/shared/errors"
	timedmanagelogic "github.com/i-Things/things/src/timed/timedjobsvr/internal/logic/timedmanage"
	"github.com/i-Things/things/src/timed/timedjobsvr/pb/timedjob"
)

type TaskSend struct {
	Code        string            `json:"code"` //任务编码
	ExecContent string            `json:"execContent"`
	Param       map[string]string `json:"param"`
}

func (s *SqlFunc) TaskSendSqlJs() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		taskMap, ok := in.Arguments[0].Export().(map[string]any)
		if !ok {
			s.Errorf("timed.SetFunc.TaskSend script use err,"+
				"need an object,code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter.AddMsg("TaskSend param not rigth"))
		}
		var task TaskSend
		err := gconv.Struct(taskMap, &task)
		if err != nil {
			s.Errorf("timed.SetFunc.TaskSend gconv.Struct err:%v",
				err)
			panic(errors.Parameter.AddMsg("TaskSend param not rigth"))
		}
		_, err = timedmanagelogic.NewTaskSendLogic(s.ctx, s.SvcCtx).TaskSend(&timedjob.TaskSendReq{
			GroupCode: s.Task.GroupCode,
			Code:      task.Code,
			ParamSql: func() *timedjob.TaskParamSql {
				if task.Param == nil {
					return nil
				}
				return &timedjob.TaskParamSql{Param: task.Param, ExecContent: task.ExecContent}
			}(),
		})
		if err != nil {
			return s.vm.ToValue(ErrRet{Err: err})
		}
		return nil
	}

}
