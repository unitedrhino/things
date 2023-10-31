package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) Hdel() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 2 {
			s.Errorf("timed.SetFunc.Hdel script use err,"+
				"need ( key, field string),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		ret, err := s.SvcCtx.Store.HdelCtx(s.ctx, s.kvKeyPre+in.Arguments[0].String(),
			in.Arguments[1].String())
		if err != nil {
			s.Errorf("timed.SetFunc.Hdel script Store.HgetCtx err:%v", err)
			panic(errors.Database.AddDetail(err))
		}
		return s.vm.ToValue(ret)
	}
}
