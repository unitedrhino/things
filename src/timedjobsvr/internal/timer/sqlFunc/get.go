package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) Get() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.Get script use err,need (key string),code:%v,script:%v",
				s.jb.Code, s.jb.Sql.Script)
			panic(errors.Parameter)
		}
		ret, err := s.SvcCtx.Store.GetCtx(s.ctx, s.kvKeyPre+in.Arguments[0].String())
		if err != nil {
			s.Errorf("timed.SetFunc.Get script Store.GetCtx err:%v", err)
			panic(err)
		}
		return s.vm.ToValue(ret)
	}
}
