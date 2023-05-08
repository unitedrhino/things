package casbin

import (
	"database/sql"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	adapter "github.com/Blank-Xu/sql-adapter"
	watcher "github.com/casbin/redis-watcher/v2"
	_ "github.com/go-sql-driver/mysql"
	r "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// casbin rule
	rule = `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
	tableName = "casbin_rule"
)

func NewCasbin(conn *sql.DB, driver string) (*casbin.Enforcer, error) {
	adapter, err := adapter.NewAdapter(conn, driver, tableName)
	logx.Must(err)

	m, err := model.NewModelFromString(rule)
	logx.Must(err)

	enforcer, err := casbin.NewEnforcer(m, adapter)
	logx.Must(err)

	err = enforcer.LoadPolicy()
	logx.Must(err)

	return enforcer, nil
}

func MustNewCasbin(conn *sql.DB, driver string) *casbin.Enforcer {
	csb, err := NewCasbin(conn, driver)
	if err != nil {
		logx.Errorw("initialize Casbin failed", logx.Field("detail", err.Error()))
		return nil
	}

	return csb
}

func NewRedisWatcher(c redis.RedisConf, f func(string2 string)) persist.Watcher {
	w, err := watcher.NewWatcher(c.Host, watcher.WatcherOptions{
		Options: r.Options{
			Network:  "tcp",
			Password: c.Pass,
		},
		Channel:    "/casbin",
		IgnoreSelf: false,
	})
	logx.Must(err)

	err = w.SetUpdateCallback(f)
	logx.Must(err)

	return w
}

func NewCasbinWithRedisWatcher(conn *sql.DB, driver string, c redis.RedisConf) *casbin.Enforcer {
	cas := MustNewCasbin(conn, driver)
	wat := NewRedisWatcher(c, func(data string) {
		watcher.DefaultUpdateCallback(cas)(data)
	})
	err := cas.SetWatcher(wat)
	logx.Must(err)
	err = cas.SavePolicy()
	logx.Must(err)
	return cas
}
