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
		ret, err := s.SvcCtx.Store.HdelCtx(s.ctx, s.GetHashKey(in.Arguments[0].String()),
			s.GetHashField(in.Arguments[1].String()))
		if err != nil {
			s.Errorf("timed.SetFunc.Hdel script Store.HgetCtx err:%v", err)
			panic(errors.Database.AddDetail(err))
		}
		//前一天的也需要删除
		ret, err = s.SvcCtx.Store.HdelCtx(s.ctx, s.GetHashKey(in.Arguments[0].String()),
			s.GetHashFieldWithDay(in.Arguments[1].String(), -1))
		if err != nil {
			s.Errorf("timed.SetFunc.Hdel script Store.HgetCtx err:%v", err)
			panic(errors.Database.AddDetail(err))
		}
		return s.vm.ToValue(ret)
	}
}
