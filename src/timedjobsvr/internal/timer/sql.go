package timer

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/task"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timedjobsvr/internal/timer/sqlFunc"
	"github.com/spf13/cast"
)

func (t Timed) SqlExec(ctx context.Context, jb *task.Info) error {
	switch jb.Sql.Type {
	case task.SqlTypeJs:
		return t.SqlJsExec(ctx, jb)
	default:
		return t.SqlNormalExec(ctx, jb)
	}
	return nil
}
func (t Timed) SqlNormalExec(ctx context.Context, jb *task.Info) error {
	dsn := cast.ToString(jb.Sql.Env[task.SqlEnvDsn])
	dbType := cast.ToString(jb.Sql.Env[task.SqlEnvDBType])
	if dsn == "" { //走默认值
		err := stores.GetCommonConn(ctx).Exec(jb.Sql.Sql).Error
		return stores.ErrFmt(err)
	}
	conn, err := stores.GetConn(conf.Database{
		DBType: dbType,
		DSN:    dsn,
	})
	if err != nil {
		return err
	}
	db, err := conn.DB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = conn.Exec(jb.Sql.Sql).Error
	return stores.ErrFmt(err)
}
func (t Timed) SqlJsExec(ctx context.Context, jb *task.Info) error {
	vm := goja.New()
	sf := sqlFunc.NewSqlFunc(ctx, t.SvcCtx, jb, vm)
	err := sf.Register()
	if err != nil {
		return err
	}
	v, err := vm.RunString(jb.Sql.Script)
	fmt.Println(v, err)
	return err
}

//func (t Timed) SetFunc(ctx context.Context, jb *task.Info, vm *goja.Runtime) error {
//	kvKeyPre := fmt.Sprintf("timed:sql:%s:%s:", jb.Group, jb.Code)
//
//	vm.Set("Set", func(in goja.FunctionCall) goja.Value {
//		if len(in.Arguments) < 2 {
//			logx.WithContext(ctx).Errorf("timed.SetFunc.Set script use err,need (key string,value string,seconds int(默认无过期时间)),code:%v,script:%v", jb.Code, jb.Sql.Script)
//			panic(errors.Parameter)
//		}
//		if len(in.Arguments) > 2 {
//			err := t.SvcCtx.Store.SetexCtx(ctx, kvKeyPre+in.Arguments[0].String(), in.Arguments[1].String(), int(in.Arguments[2].ToInteger()))
//			if err != nil {
//				logx.WithContext(ctx).Errorf("timed.SetFunc.Set script Store.GetCtx err:%v", err)
//				panic(err)
//			}
//		}
//		err := t.SvcCtx.Store.SetCtx(ctx, kvKeyPre+in.Arguments[0].String(), in.Arguments[1].String())
//		if err != nil {
//			logx.WithContext(ctx).Errorf("timed.SetFunc.Set script Store.GetCtx err:%v", err)
//			panic(err)
//		}
//		return nil
//	})
//	vm.Set("Get", func(in goja.FunctionCall) goja.Value {
//		if len(in.Arguments) < 1 {
//			logx.WithContext(ctx).Errorf("timed.SetFunc.Get script use err,need (key string),code:%v,script:%v", jb.Code, jb.Sql.Script)
//			panic(errors.Parameter)
//		}
//		ret, err := t.SvcCtx.Store.GetCtx(ctx, kvKeyPre+in.Arguments[0].String())
//		if err != nil {
//			logx.WithContext(ctx).Errorf("timed.SetFunc.Get script Store.GetCtx err:%v", err)
//			panic(err)
//		}
//		return vm.ToValue(ret)
//	})
//	vm.Set("Log", func(in goja.FunctionCall) goja.Value {
//		logx.WithContext(ctx).Infof("script  code:%v log:%v", jb.Code, in.Arguments)
//		return nil
//	})
//	vm.Set("GetEnv", func(in goja.FunctionCall) goja.Value {
//		if len(in.Arguments) < 1 {
//			logx.WithContext(ctx).Errorf("timed.SetFunc.GetEnv script use err,need (key string),code:%v,script:%v", jb.Code, jb.Sql.Script)
//			panic(errors.Parameter)
//		}
//		return vm.ToValue(jb.Sql.Env[in.Arguments[0].String()])
//	})
//	getConn := func(in goja.FunctionCall) (*gorm.DB, func() error) {
//		dsn := cast.ToString(jb.Sql.Env[task.SqlEnvDsn])
//		dbType := cast.ToString(jb.Sql.Env[task.SqlEnvDBType])
//		if len(in.Arguments) > 1 {
//			dsn = in.Arguments[1].String()
//			if len(in.Arguments) > 2 {
//				dbType = in.Arguments[2].String()
//			}
//		}
//		if dsn == "" { //走默认值
//			return stores.GetCommonConn(ctx), func() error {
//				return nil
//			}
//		}
//		conn, err := stores.GetConn(conf.Database{
//			DBType: dbType,
//			DSN:    dsn,
//		})
//		if err != nil {
//			panic(err)
//		}
//		db, err := conn.DB()
//		if err != nil {
//			panic(err)
//		}
//		return conn, db.Close
//	}
//	//第一个参数是sql 第二个参数是dsn(可选),第三个参数是dbType(默认mysql)
//	vm.Set("Exec", func(in goja.FunctionCall) goja.Value {
//		if len(in.Arguments) < 1 {
//			logx.WithContext(ctx).Errorf("timed.SetFunc.Exec script use err,need (第一个参数是sql 第二个参数是dsn(可选),第三个参数是dbType(默认mysql)),code:%v,script:%v", jb.Code, jb.Sql.Script)
//			panic(errors.Parameter)
//		}
//		sql := in.Arguments[0].String()
//		conn, close := getConn(in)
//		defer close()
//		err := conn.Exec(sql).Error
//		if err != nil {
//			panic(err)
//		}
//		return nil
//	})
//	//第一个参数是sql 第二个参数是dsn(可选),第三个参数是dbType(默认mysql)
//	vm.Set("Select", func(in goja.FunctionCall) goja.Value {
//		if len(in.Arguments) < 1 {
//			logx.WithContext(ctx).Errorf("timed.SetFunc.Select script use err,need (第一个参数是sql 第二个参数是dsn(可选),第三个参数是dbType(默认mysql)),code:%v,script:%v", jb.Code, jb.Sql.Script)
//			panic(errors.Parameter)
//		}
//		sql := in.Arguments[0].String()
//		conn, close := getConn(in)
//		defer close()
//		var ret []map[string]any
//		err := conn.Raw(sql).Scan(&ret).Error
//		if err != nil {
//			panic(err)
//		}
//		return vm.ToValue(ret)
//	})
//	return nil
//}
