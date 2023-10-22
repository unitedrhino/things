package sqlFunc

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/i-Things/things/shared/conf"
	"github.com/i-Things/things/shared/domain/task"
	"github.com/i-Things/things/shared/stores"
	"github.com/i-Things/things/src/timedjobsvr/internal/domain"
	"github.com/i-Things/things/src/timedjobsvr/internal/svc"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type SqlFunc struct {
	SvcCtx     *svc.ServiceContext
	ctx        context.Context
	jb         *task.Info
	vm         *goja.Runtime
	ExecuteLog []*domain.ScriptLog
	kvKeyPre   string
	logx.Logger
}

func NewSqlFunc(ctx context.Context, svcCtx *svc.ServiceContext, jb *task.Info, vm *goja.Runtime) *SqlFunc {
	kvKeyPre := fmt.Sprintf("timed:sql:%s:%s:", jb.Group, jb.Code)
	if code := jb.Sql.Env["code"]; code != nil {
		kvKeyPre = fmt.Sprintf("timed:sql:%s:%s:", jb.Group, cast.ToString(code))
	}
	return &SqlFunc{SvcCtx: svcCtx, ctx: ctx, Logger: logx.WithContext(ctx), jb: jb, vm: vm, kvKeyPre: kvKeyPre}
}

func (s *SqlFunc) Register() error {
	var funcList = []struct {
		Name string
		f    func(in goja.FunctionCall) goja.Value
	}{
		{"Set", s.Set()},
		{"Get", s.Get()},
		{"Select", s.Select()},
		{"Exec", s.Exec()},
		{"LogError", s.LogError()},
		{"LogInfo", s.LogInfo()},
		{"GetEnv", s.GetEnv()},
		{"Hexists", s.Hexists()},
		{"Hdel", s.Hdel()},
		{"HGet", s.HGet()},
		{"HSet", s.HSet()},
		{"HGetAll", s.HGetAll()},
	}
	for _, f := range funcList {
		err := s.vm.Set(f.Name, f.f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SqlFunc) getConn(in goja.FunctionCall) (*gorm.DB, func() error) {
	dsn := cast.ToString(s.jb.Sql.Env[task.SqlEnvDsn])
	dbType := cast.ToString(s.jb.Sql.Env[task.SqlEnvDBType])
	if len(in.Arguments) > 1 {
		dsn = in.Arguments[1].String()
		if len(in.Arguments) > 2 {
			dbType = in.Arguments[2].String()
		}
	}
	if dsn == "" { //走默认值
		return stores.GetCommonConn(s.ctx), func() error {
			return nil
		}
	}
	conn, err := stores.GetConn(conf.Database{
		DBType: dbType,
		DSN:    dsn,
	})
	if err != nil {
		panic(err)
	}
	db, err := conn.DB()
	if err != nil {
		panic(err)
	}
	return conn, db.Close
}
