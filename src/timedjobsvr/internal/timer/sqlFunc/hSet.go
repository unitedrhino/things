package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) HSet() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 3 {
			s.Errorf("timed.SetFunc.HSet script use err,"+
				"need (key, field, value string),code:%v,script:%v",
				s.jb.Code, s.jb.Sql.Script)
			panic(errors.Parameter)
		}
		err := s.SvcCtx.Store.HsetCtx(s.ctx, s.kvKeyPre+in.Arguments[0].String(),
			in.Arguments[1].String(), in.Arguments[2].String())
		if err != nil {
			s.Errorf("timed.SetFunc.Set script Store.HsetCtx err:%v", err)
			panic(err)
		}
		return nil
	}

}
