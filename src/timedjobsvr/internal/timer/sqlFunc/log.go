package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/utils"
)

func (s *SqlFunc) Log() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		var args []any
		for _, arg := range in.Arguments {
			args = append(args, arg.Export())
		}
		s.Infof("script  code:%v log:%v", s.jb.Code, utils.Fmt(args))
		return nil
	}

}
