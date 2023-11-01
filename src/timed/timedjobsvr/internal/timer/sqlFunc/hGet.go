package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
	"strings"
)

func (s *SqlFunc) HGet() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 2 {
			s.Errorf("timed.SetFunc.HGet script use err,"+
				"need (key, field string),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		ret, err := s.SvcCtx.Store.HgetCtx(s.ctx, s.GetHashKey(in.Arguments[0].String()),
			s.GetHashField(in.Arguments[1].String()))
		if err != nil {
			if strings.Contains(err.Error(), "redis: nil") {
				ret, err = s.SvcCtx.Store.HgetCtx(s.ctx, s.GetHashKey(in.Arguments[0].String()),
					s.GetHashFieldWithDay(in.Arguments[1].String(), -1))
			}
			if err != nil && !strings.Contains(err.Error(), "redis: nil") {
				s.Errorf("timed.SetFunc.Set script Store.HgetCtx err:%v", err)
				panic(errors.Database.AddDetail(err))
			}
		}
		return s.vm.ToValue(ret)
	}

}
