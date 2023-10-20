package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) HGetAll() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.HGetAll script use err,"+
				"need (key string),code:%v,script:%v",
				s.jb.Code, s.jb.Sql.Script)
			panic(errors.Parameter)
		}
		ret, err := s.SvcCtx.Store.HgetallCtx(s.ctx, s.kvKeyPre+in.Arguments[0].String())
		if err != nil {
			s.Errorf("timed.SetFunc.Set script Store.HgetCtx err:%v", err)
			panic(err)
		}
		return s.vm.ToValue(ret)
	}

}
