package timedDelay

import (
	"encoding/json"
	"fmt"
	"github.com/i-Things/things/shared/domain/task"
)

type QueueMsg struct {
	Type    string //消息队列的类型 natsJs nats
	Topic   string //发送的主题
	Payload any    //发送的消息内容
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
	err := t.Enqueue(&j, option)
	return err
}
