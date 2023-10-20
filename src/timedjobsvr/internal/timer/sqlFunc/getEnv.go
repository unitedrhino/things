package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) GetEnv() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.GetEnv script use err,need (key string),code:%v,script:%v",
				s.jb.Code, s.jb.Sql.Script)
			panic(errors.Parameter)
		}
		return s.vm.ToValue(s.jb.Sql.Env[in.Arguments[0].String()])
	}
}
