package timedDelay

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/domain/task"
)

type QueueMsg struct {
	Type     string //消息队列的类型 natsJs nats
	Priority string //优先级: 6:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	Topic    string //发送的主题
	Payload  any    //发送的消息内容
}

// 延时发送消息
func (t Timed) DelayQueue(msg QueueMsg, option *Option) error {
	payload, _ := json.Marshal(msg.Payload)
	switch msg.Payload.(type) {
	case string:
		payload = []byte(msg.Payload.(string))
	case []byte:
		payload = msg.Payload.([]byte)
	}
	params, _ := json.Marshal(task.Queue{
		Topic:   msg.Topic,
		Payload: string(payload),
	})
	j := task.Info{
		Group:   t.serverName,
		Type:    task.TaskTypeQueue,
		Code:    fmt.Sprintf("delayQueue_%s_%s", t.serverName, msg.Topic),
		SubType: msg.Type,
		Name:    "服务延时消息推送",
		Params:  string(params),
	}
	err := j.Init()
	if err != nil {
		return err
	}
	var opts []asynq.Option
	if option != nil {
		var opt = asynq.ProcessAt(option.ProcessAt)
		if option.ProcessIn != 0 {
			opt = asynq.ProcessIn(option.ProcessIn)
		}
		opts = append(opts, opt)
		if option.Timeout != 0 {
			opts = append(opts, asynq.Timeout(option.Timeout))
		}
		if !option.Deadline.IsZero() {
			opts = append(opts, asynq.Deadline(option.Deadline))
		}
	}

	_, err = t.asynqClient.Enqueue(j.ToTask(), opts...)
	return err
}
