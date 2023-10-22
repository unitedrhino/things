package timer

import (
	"context"
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/task"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timedjobsvr/internal/repo/relationDB"
	"github.com/i-Things/things/src/timedjobsvr/internal/timer/sqlFunc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
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
	err := func() error {
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
	}()
	e := errors.Fmt(err)
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedJobLog{
		Group:      jb.Group,
		Type:       jb.Type,
		SubType:    jb.SubType,
		Name:       jb.Name,
		Code:       jb.Code,
		ResultCode: e.GetCode(),
		ResultMsg:  e.GetMsg(),
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("SqlNormalExec.JobLog.Insert err:%v", er)
	}
	return err
}

func (t Timed) SqlJsExec(ctx context.Context, jb *task.Info) error {
	var SqlJob func() map[string]any
	vm := goja.New()
	sf := sqlFunc.NewSqlFunc(ctx, t.SvcCtx, jb, vm)
	err := func() error {
		err := sf.Register()
		if err != nil {
			return err
		}
		_, err = vm.RunString(jb.Sql.Script)
		if err != nil {
			return err
		}
		err = vm.ExportTo(vm.Get("SqlJob"), &SqlJob)
		if err != nil {
			return err
		}
		return nil
	}()
	var (
		code = errors.OK.Code
		msg  = errors.OK.Msg
	)
	e := errors.Fmt(err)
	if e != nil {
		code = e.GetCode()
		msg = e.GetMsg()
	} else if SqlJob != nil {
		func() {
			defer func() {
				if p := recover(); p != nil {
					if e, ok := p.(*goja.Exception); !ok { //如果不是sql执行错误
						panic(p)
					} else {
						er := e.Error()
						code = errors.Script.Code
						msg = errors.Script.AddMsg(er).GetMsg()
					}
				}
			}()
			ret := SqlJob()
			code = cast.ToInt64(ret["code"])
			msg = cast.ToString(ret["msg"])
		}()
	}

	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedJobLog{
		Group:      jb.Group,
		Type:       jb.Type,
		SubType:    jb.SubType,
		Name:       jb.Name,
		Code:       jb.Code,
		ResultCode: code,
		ResultMsg:  msg,
		ExecLog:    sf.ExecuteLog,
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("SqlNormalExec.JobLog.Insert err:%v", er)
	}
	return er
}
