package sqlFunc

import (
	"encoding/json"
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
	"github.com/spf13/cast"
)

func (s *SqlFunc) Hset() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 3 {
			s.Errorf("timed.SetFunc.Hset script use err,"+
				"need (key, field, value string),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		v := in.Arguments[2].Export()
		value, err := cast.ToStringE(v)
		if err != nil {
			b, err := json.Marshal(v)
			if err != nil {
				return s.vm.ToValue(err)
			}
			value = string(b)
		}
		err = s.SvcCtx.Store.HsetCtx(s.ctx, s.GetHashKey(in.Arguments[0].String()),
			s.GetHashField(in.Arguments[1].String()), value)
		if err != nil {
			s.Errorf("timed.SetFunc.Set script Store.HsetCtx err:%v", err)
			return s.vm.ToValue(err)
		}
		return nil
	}

}
