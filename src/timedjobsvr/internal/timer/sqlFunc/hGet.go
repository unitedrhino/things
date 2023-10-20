package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) HGet() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 2 {
			s.Errorf("timed.SetFunc.HGet script use err,"+
				"need (key, field string),code:%v,script:%v",
				s.jb.Code, s.jb.Sql.Script)
			panic(errors.Parameter)
		}
		ret,err := s.SvcCtx.Store.HgetCtx(s.ctx, s.kvKeyPre+in.Arguments[0].String(),
			in.Arguments[1].String())
		if err != nil {
			s.Errorf("timed.SetFunc.Set script Store.HgetCtx err:%v", err)
			panic(err)
		}
		return s.vm.ToValue(ret)
	}

}
