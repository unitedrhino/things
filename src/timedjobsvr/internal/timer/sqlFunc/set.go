package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) Set() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 2 {
			s.Errorf("timed.SetFunc.Set script use err,"+
				"need (key string,value string,seconds int(默认无过期时间)),code:%v,script:%v",
				s.jb.Code, s.jb.Sql.Script)
			panic(errors.Parameter)
		}
		if len(in.Arguments) > 2 {
			err := s.SvcCtx.Store.SetexCtx(s.ctx, s.kvKeyPre+in.Arguments[0].String(),
				in.Arguments[1].String(), int(in.Arguments[2].ToInteger()))
			if err != nil {
				s.Errorf("timed.SetFunc.Set script Store.GetCtx err:%v", err)
				panic(err)
			}
		}
		err := s.SvcCtx.Store.SetCtx(s.ctx, s.kvKeyPre+in.Arguments[0].String(), in.Arguments[1].String())
		if err != nil {
			s.Errorf("timed.SetFunc.Set script Store.GetCtx err:%v", err)
			panic(err)
		}
		return nil
	}

}
