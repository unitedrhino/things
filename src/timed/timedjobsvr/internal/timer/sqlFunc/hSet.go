package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) Hset() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 3 {
			s.Errorf("timed.SetFunc.Hset script use err,"+
				"need (key, field, value string),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		err := s.SvcCtx.Store.HsetCtx(s.ctx, s.GetHashKey(in.Arguments[0].String()),
			s.GetHashField(in.Arguments[1].String()), in.Arguments[2].String())
		if err != nil {
			s.Errorf("timed.SetFunc.Set script Store.HsetCtx err:%v", err)
			panic(errors.Database.AddDetail(err))
		}
		return nil
	}

}
