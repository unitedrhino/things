package timer

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/errors"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timed/internal/domain"
	"github.com/i-Things/things/src/timed/internal/repo/relationDB"
	"github.com/i-Things/things/src/timed/timedjobsvr/internal/timer/sqlFunc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

func (t Timed) SqlExec(ctx context.Context, task *domain.TaskInfo) error {
	switch task.GroupSubType {
	case domain.SqlTypeJs:
		return t.SqlJsExec(ctx, task)
	default:
		return t.SqlNormalExec(ctx, task)
	}
	return nil
}
func (t Timed) SqlNormalExec(ctx context.Context, task *domain.TaskInfo) error {
	err := func() error {
		dsn := cast.ToString(task.Sql.Env[domain.SqlEnvDsn])
		dbType := cast.ToString(task.Sql.Env[domain.SqlEnvDBType])
		if dsn == "" { //走默认值
			err := stores.GetCommonConn(ctx).Exec(task.Sql.Param.ExecContent).Error
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
		err = conn.Exec(task.Sql.Param.ExecContent).Error
		return stores.ErrFmt(err)
	}()
	e := errors.Fmt(err)
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedTaskLog{
		ResultCode: e.GetCode(),
		ResultMsg:  e.GetMsg(),
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("SqlNormalExec.JobLog.Insert err:%v", er)
	}
	return err
}

func (t Timed) SqlJsExec(ctx context.Context, task *domain.TaskInfo) error {
	var (
		code = errors.OK.Code
		msg  = errors.OK.Msg
	)
	vm := goja.New()
	sf := sqlFunc.NewSqlFunc(ctx, t.SvcCtx, task, vm)

	func() {
		defer func() {
			if p := recover(); p != nil {
				if e, ok := p.(error); ok {
					err := errors.Fmt(e)
					sf.ExecuteLog = append(sf.ExecuteLog, &domain.ScriptLog{
						Level:       "error",
						Content:     fmt.Sprintf("catch an panic,err:%v", err.Error()),
						CreatedTime: time.Now().Unix(),
					})
					code = err.GetCode()
					msg = err.GetMsg()
				}
			}
		}()
		var SqlJob func() map[string]any
		err := func() error {
			err := sf.Register()
			if err != nil {
				return err
			}
			_, err = vm.RunString(task.Sql.Param.ExecContent)
			if err != nil {
				return err
			}
			err = vm.ExportTo(vm.Get("SqlJob"), &SqlJob)
			if err != nil {
				return err
			}
			return nil
		}()

		e := errors.Fmt(err)
		if e != nil {
			code = e.GetCode()
			msg = e.GetMsg()
		} else if SqlJob != nil {
			ret := SqlJob()
			code = cast.ToInt64(ret["code"])
			msg = cast.ToString(ret["msg"])
		}
	}()
	er := relationDB.NewJobLogRepo(ctx).Insert(ctx, &relationDB.TimedTaskLog{
		GroupCode:  task.GroupCode,
		TaskCode:   task.Code,
		ResultCode: code,
		ResultMsg:  msg,
		TimedTaskSqlScript: &relationDB.TimedTaskSqlScript{
			ExecLog:   sf.ExecuteLog,
			SelectNum: sf.SelectNum,
			ExecNum:   sf.ExecNum,
		},
	})
	if er != nil {
		logx.WithContext(ctx).Errorf("SqlNormalExec.JobLog.Insert err:%v", er)
	}
	return er
}
