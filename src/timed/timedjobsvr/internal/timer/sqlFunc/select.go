package sqlFunc

import (
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/errors"
)

func (s *SqlFunc) Select() func(in goja.FunctionCall) goja.Value {
	return func(in goja.FunctionCall) goja.Value {
		if len(in.Arguments) < 1 {
			s.Errorf("timed.SetFunc.Select script use err,"+
				"need (第一个参数是sql 第二个参数是指定的数据库连接code(可选,不填选择默认的连接,需要在config里配置),code:%v,script:%v",
				s.Task.Code, s.Task.Sql.Param.ExecContent)
			panic(errors.Parameter)
		}
		sql := in.Arguments[0].String()
		conn, close := s.getConn(in, "select")
		defer close()
		var ret []map[string]any
		v := conn.Raw(sql).Scan(&ret)
		err := v.Error
		s.SelectNum += v.RowsAffected
		if err != nil {
			panic(errors.Database.AddDetail(err))
		}
		return s.vm.ToValue(ret)
	}

}
