package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/utils"
	"github.com/i-Things/things/src/timed/internal/domain"
	"time"
)

func (s *SqlFunc) LogError() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		var args []any
		for _, arg := range in.Arguments {
			args = append(args, arg.Export())
		}
		s.ExecuteLog = append(s.ExecuteLog, &domain.ScriptLog{
			Level:       "error",
			Content:     utils.Fmt(args),
			CreatedTime: time.Now().Unix(),
		})
		s.Errorf("script  code:%v log:%v", s.Task.Code, utils.Fmt(args))
		return nil
	}
}

func (s *SqlFunc) LogInfo() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		var args []any
		for _, arg := range in.Arguments {
			args = append(args, arg.Export())
		}
		s.ExecuteLog = append(s.ExecuteLog, &domain.ScriptLog{
			Level:       "info",
			Content:     utils.Fmt(args),
			CreatedTime: time.Now().Unix(),
		})
		s.Infof("script  code:%v log:%v", s.Task.Code, utils.Fmt(args))
		return nil
	}
}
