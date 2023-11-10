package task

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/errors"
)

const (
	TaskTypeQueue = "queue"
	TaskTypeSql   = "sql"
)

type Info struct {
	Group    string `json:"group"`    // 任务组名
	Type     string `json:"type"`     //任务类型:queue(消息队列消息发送)  sql(执行sql) email(邮件发送) http(http请求)
	SubType  string `json:"subType"`  //任务子类型 natsJs nats
	Name     string `json:"name"`     // 任务名称
	Code     string `json:"code"`     //任务编码
	Params   string `json:"params"`   // 任务参数
	Priority string `json:"priority"` //优先级: 6:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	Queue    *Queue `json:"-"`        //消息队列类型
	Sql      *Sql   `json:"-"`        //sql执行类型
}

func (j *Info) GetTypeName() string {
	return fmt.Sprintf("%s:%s", j.Group, j.Code)
}

func (j *Info) ToPayload() []byte {
	ret, _ := json.Marshal(j)
	return ret
}
func (j *Info) ToPriority() string {
	if j.Priority == "" {
		return "default"
	}
	return j.Priority
}

func (j *Info) Init() error {
	switch j.Type {
	case TaskTypeQueue:
		var q Queue
		err := json.Unmarshal([]byte(j.Params), &q)
		if err != nil {
			return err
		}
		j.Queue = &q
		return nil
	case TaskTypeSql:
		var s Sql
		err := json.Unmarshal([]byte(j.Params), &s)
		if err != nil {
			return err
		}
		j.Sql = &s
		return nil
	}
	return errors.Parameter.AddMsgf("job type not support:%v", j.Type)
}

func (j Info) ToTask() *asynq.Task {
	return asynq.NewTask(j.GetTypeName(), j.ToPayload(), asynq.Queue(j.ToPriority()))
}
