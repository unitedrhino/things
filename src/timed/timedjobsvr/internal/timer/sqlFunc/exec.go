package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
)

func (s *SqlFunc) Exec() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.Exec script use err,"+
				"need (第一个参数是sql 第二个参数是指定的数据库连接code(可选,不填选择默认的连接,需要在config里配置),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		sql := in.Arguments[0].String()
		conn, close := s.getConn(in, "exec")
		defer close()
		ret := conn.Exec(sql)
		err := ret.Error
		s.ExecNum += ret.RowsAffected
		if err != nil {
			return s.vm.ToValue(ErrRet{Err: stores.ErrFmt(err)})
		}
		return s.vm.ToValue(ErrRet{})
	}

}
