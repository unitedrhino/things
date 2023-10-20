package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) Exec() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.Exec script use err,"+
				"need (第一个参数是sql 第二个参数是dsn(可选),第三个参数是dbType(默认mysql)),code:%v,script:%v",
				s.jb.Code, s.jb.Sql.Script)
			panic(errors.Parameter)
		}
		sql := in.Arguments[0].String()
		conn, close := s.getConn(in)
		defer close()
		err := conn.Exec(sql).Error
		if err != nil {
			panic(err)
		}
		return nil
	}

}
