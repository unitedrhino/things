package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) Create() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 2 {
			s.Errorf("timed.SetFunc.Exec script use err,"+
				"need (第一个参数是sql 第二个参数是指定的数据库连接code(可选,不填选择默认的连接,需要在config里配置),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		table := in.Arguments[0].String()
		data := in.Arguments[1].Export()
		conn, close := s.getConn(in, "exec")
		defer close()
		ret := conn.Table(table).Create(data)
		err := ret.Error
		s.ExecNum += ret.RowsAffected
		if err != nil {
			return s.vm.ToValue(err)
		}
		return nil
	}

}
