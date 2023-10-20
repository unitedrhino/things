package timedDelay

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/domain/task"
)

type SqlMsg struct {
	Type   string // sql执行的类型  normal(支持环境变量的直接执行方式),JavaScript(js脚本执行)
	Sql    string // 执行的sql
	Script string `json:"script"`
	//需要使用的环境变量 dsn:如果填写,默认使用该dsn来连接,dbType:mysql|pgsql 默认是mysql
	Env map[string]any `json:"env"`
}

func (t Timed) SqlExec(msg task.Sql, option *Option) error {
	params, _ := json.Marshal(task.Sql{
		Sql:  msg.Sql,
		Type: msg.Type,
		Env:  msg.Env,
	})
	j := task.Info{
		Group:   t.serverName,
		Type:    task.TaskTypeQueue,
		Code:    fmt.Sprintf("delaySql_%s_%s", t.serverName, msg.Type),
		SubType: msg.Type,
		Name:    "服务延时消息推送",
		Params:  string(params),
	}
	err := t.Enqueue(&j, option)
	return err
}
