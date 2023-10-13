package timedDelay

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/i-Things/things/shared/domain/task"
	"time"
)

type QueueMsg struct {
	Type     string //消息队列的类型 natsJs nats
	Priority string //优先级: 6:critical 最高优先级  3: default 普通优先级 1:low 低优先级
	Topic    string //发送的主题
	Payload  any    //发送的消息内容
	//以下两个参数优先使用ProcessIn
	ProcessIn time.Duration //多久之后发
	ProcessAt time.Time     // 固定时间发
}

// 延时发送消息
func (t Timed) DelayQueue(msg QueueMsg) error {
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
	var opt = asynq.ProcessAt(msg.ProcessAt)
	if msg.ProcessIn != 0 {
		opt = asynq.ProcessIn(msg.ProcessIn)
	}
	_, err = t.asynqClient.Enqueue(j.ToTask(), opt)
	return err
}
